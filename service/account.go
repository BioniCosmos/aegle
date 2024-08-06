package service

import (
	"bytes"
	"context"
	cryptoRand "crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"image/jpeg"
	mathRand "math/rand"
	"net/smtp"
	"time"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/setting"
	"github.com/bionicosmos/argon2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrEmailExists         = fiber.NewError(fiber.StatusConflict, "Email exists")
	ErrAccountNotExist     = fiber.NewError(fiber.StatusNotFound, "Account does not exist")
	ErrPassword            = fiber.NewError(fiber.StatusBadRequest, "Password mismatch")
	ErrVerified            = fiber.NewError(fiber.StatusConflict, "Verified account")
	ErrLinkExpired         = fiber.NewError(fiber.StatusNotFound, "Link expired")
	ErrTOTPNotFound        = fiber.NewError(fiber.StatusBadRequest, "Uninitialized or expired TOTP")
	ErrInvalidTOTP         = fiber.NewError(fiber.StatusUnauthorized, "Invalid TOTP")
	ErrDisabledAuthType    = fiber.NewError(fiber.StatusConflict, "The authentication type not enabled")
	ErrInvalidMFACode      = fiber.NewError(fiber.StatusUnauthorized, "Invalid multi-factor authentication code")
	ErrNoNeedRecoveryCode  = fiber.NewError(fiber.StatusConflict, "No need to generate recovery codes")
	ErrInvalidRecoveryCode = fiber.NewError(fiber.StatusBadRequest, "Invalid recovery code")
)

var tmpl *template.Template

func SignUp(body *transfer.SignUpBody) (model.Account, error) {
	if _, err := model.FindAccount(body.Email); err != nil &&
		!errors.Is(err, mongo.ErrNoDocuments) {
		return model.Account{}, err
	} else if err == nil {
		return model.Account{}, ErrEmailExists
	}
	return transactionWithValue(
		func(ctx mongo.SessionContext) (model.Account, error) {
			account := model.Account{
				Email:    body.Email,
				Name:     body.Name,
				Password: argon2.Hash(body.Password),
				Role:     model.AccountRoleMember,
				Status:   model.AccountStatusUnverified,
			}
			if err := model.InsertAccount(ctx, &account); err != nil {
				return model.Account{}, err
			}
			id := uuid.NewString()
			link := model.VerificationLink{
				Id:        id,
				Email:     body.Email,
				CreatedAt: time.Now(),
			}
			if err := model.InsertVerificationLink(ctx, &link); err != nil {
				return model.Account{}, err
			}
			if err := sendLinkEmail(&link); err != nil {
				return model.Account{}, err
			}
			return account, nil
		},
	)
}

func Verify(id string, email string) error {
	a, err := model.FindAccount(email)
	if err != nil {
		return err
	}
	if err := checkStatus(&a); err != nil {
		return err
	}
	if err := model.FindVerificationLink(id); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrLinkExpired
		}
		return err
	}
	return transaction(func(ctx mongo.SessionContext) error {
		if err := model.DeleteVerificationLinks(ctx, a.Email); err != nil {
			return err
		}
		return model.UpdateAccount(
			ctx,
			bson.M{"email": a.Email},
			bson.M{"$set": bson.M{"status": model.AccountStatusNormal}},
		)
	})
}

func SendVerificationLink(email string) error {
	account, err := model.FindAccount(email)
	if err != nil {
		return err
	}
	if err := checkStatus(&account); err != nil {
		return err
	}
	return transaction(func(ctx mongo.SessionContext) error {
		if err := model.DeleteVerificationLinks(
			ctx,
			account.Email,
		); err != nil {
			return err
		}
		link := model.VerificationLink{
			Id:        uuid.NewString(),
			Email:     account.Email,
			CreatedAt: time.Now(),
		}
		if err := model.InsertVerificationLink(ctx, &link); err != nil {
			return err
		}
		return sendLinkEmail(&link)
	})
}

func SignIn(body *transfer.SignInBody) (model.Account, error) {
	account, err := model.FindAccount(body.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Account{}, ErrAccountNotExist
		}
		return model.Account{}, err
	}
	if !argon2.Verify(body.Password, account.Password) {
		return model.Account{}, ErrPassword
	}
	return account, nil
}

func MFA(email string, mfaBody *transfer.MFABody) error {
	account, err := model.FindAccount(email)
	if err != nil {
		return err
	}
	switch mfaBody.Type {
	case transfer.MFATypeTOTP:
		if account.MFA.TOTP == "" {
			return ErrDisabledAuthType
		}
		if !totp.Validate(mfaBody.Code, account.MFA.TOTP) {
			return ErrInvalidMFACode
		}
	case transfer.MFATypeRecoveryCode:
		if !account.MFA.RecoveryCode {
			return ErrDisabledAuthType
		}
		filter := bson.M{"email": email, "code": mfaBody.Code}
		code, err := model.FindRecoveryCode(context.Background(), filter)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return ErrInvalidRecoveryCode
			}
			return err
		}
		if code.Used {
			return ErrInvalidRecoveryCode
		}
		if err := model.UpdateRecoveryCode(
			context.Background(),
			filter,
			bson.M{"$set": bson.M{"used": true}},
		); err != nil {
			return err
		}
	default:
		return ErrDisabledAuthType
	}
	return nil
}

func CreateTOTP(email string) (transfer.CreateTOTPBody, error) {
	return transactionWithValue(
		func(ctx mongo.SessionContext) (transfer.CreateTOTPBody, error) {
			if err := model.DeleteTOTP(ctx, email); err != nil {
				return transfer.CreateTOTPBody{}, err
			}
			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "Aegle",
				AccountName: email,
			})
			if err != nil {
				return transfer.CreateTOTPBody{}, err
			}
			if err := model.InsertTOTP(
				ctx,
				&model.TOTP{
					Email:     email,
					Secret:    key.Secret(),
					CreatedAt: time.Now(),
				},
			); err != nil {
				return transfer.CreateTOTPBody{}, err
			}
			image, err := key.Image(200, 200)
			if err != nil {
				return transfer.CreateTOTPBody{}, err
			}
			buf := bytes.Buffer{}
			if err := jpeg.Encode(&buf, image, nil); err != nil {
				return transfer.CreateTOTPBody{}, err
			}
			base64.StdEncoding.EncodeToString(buf.Bytes())
			return transfer.CreateTOTPBody{
				Secret: key.Secret(),
				Image: "data:image/jpeg;base64," +
					base64.StdEncoding.EncodeToString(buf.Bytes()),
			}, nil
		},
	)
}

func ConfirmTOTP(email string, code string) error {
	t, err := model.FindTOTP(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrTOTPNotFound
		}
		return err
	}
	if !totp.Validate(code, t.Secret) {
		return ErrInvalidTOTP
	}
	return transaction(func(ctx mongo.SessionContext) error {
		if err := model.DeleteTOTP(ctx, email); err != nil {
			return err
		}
		return model.UpdateAccount(
			ctx,
			bson.M{"email": email},
			bson.M{"$set": bson.M{"mfa": bson.M{"totp": t.Secret}}},
		)
	})
}

func DeleteTOTP(email string) error {
	return transaction(func(ctx mongo.SessionContext) error {
		if err := model.DeleteTOTP(ctx, email); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return ErrTOTPNotFound
			}
			return err
		}
		return model.UpdateAccount(
			ctx,
			bson.M{"email": email},
			bson.M{"$unset": bson.M{"mfa": bson.M{"totp": ""}}},
		)
	})
}

func GetRecoveryCodes(email string) ([]model.RecoveryCode, error) {
	account, err := model.FindAccount(email)
	if err != nil {
		return nil, err
	}
	if account.MFA.TOTP == "" {
		return nil, ErrNoNeedRecoveryCode
	}
	count, err := model.CountRecoveryCodes(context.Background(), email)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return transactionWithValue(
			func(ctx mongo.SessionContext) ([]model.RecoveryCode, error) {
				if err := model.UpdateAccount(
					ctx,
					bson.M{"email": email},
					bson.M{"$set": bson.M{"mfa.recoveryCode": true}},
				); err != nil {
					return nil, err
				}
				return RenewRecoveryCodes(email)
			},
		)
	}
	return model.FindRecoveryCodes(context.Background(), email)
}

func RenewRecoveryCodes(email string) ([]model.RecoveryCode, error) {
	account, err := model.FindAccount(email)
	if err != nil {
		return nil, err
	}
	if account.MFA.TOTP == "" {
		return nil, ErrNoNeedRecoveryCode
	}
	codes := []model.RecoveryCode{}
	newCodes, err := generateRecoveryCodes(10, 8)
	if err != nil {
		return nil, err
	}
	for _, code := range newCodes {
		codes = append(
			codes,
			model.RecoveryCode{Code: code, Used: false, Email: email},
		)
	}
	return codes, transaction(func(ctx mongo.SessionContext) error {
		if err := model.DeleteRecoveryCodes(ctx, email); err != nil {
			return err
		}
		return model.InsertRecoveryCodes(ctx, codes)
	})
}

//lint:ignore U1000 todo
func generateCode(length int) string {
	digits := make([]byte, length)
	for i := range length {
		digits[i] = byte(mathRand.Intn(10)) + '0'
	}
	return string(digits)
}

func checkStatus(a *model.Account) error {
	if a.Status != model.AccountStatusUnverified {
		return ErrVerified
	}
	return nil
}

func sendLinkEmail(link *model.VerificationLink) error {
	body := bytes.Buffer{}
	if err := tmpl.Execute(
		&body,
		setting.X.BaseURL+"/dashboard/verification/"+link.Id,
	); err != nil {
		return err
	}
	return sendMail(link.Email, "Sign-up Verification", &body)
}

func sendMail(to string, subject string, body *bytes.Buffer) error {
	email := setting.X.Email
	auth := smtp.PlainAuth("", email.Username, email.Password, email.Host)
	message := bytes.Buffer{}
	message.WriteString("Subject: ")
	message.WriteString(subject)
	message.WriteString("\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
	body.WriteTo(&message)
	return smtp.SendMail(
		fmt.Sprintf("%v:%v", email.Host, email.Port),
		auth,
		email.Username,
		[]string{to},
		message.Bytes(),
	)
}

func generateRecoveryCode(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := cryptoRand.Read(bytes)
	if err != nil {
		return "", err
	}
	encoded := base32.
		StdEncoding.
		WithPadding(base32.NoPadding).
		EncodeToString(bytes)
	return encoded[:n], nil
}

func generateRecoveryCodes(count, length int) ([]string, error) {
	codes := []string{}
	for i := 0; i < count; i++ {
		code, err := generateRecoveryCode(length)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}
