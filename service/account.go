package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"image/jpeg"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/model/account"
	"github.com/bionicosmos/aegle/setting"
	"github.com/bionicosmos/argon2"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrAccountExists = errors.New("account exists")
	ErrPassword      = errors.New("password mismatch")
	ErrVerified      = errors.New("verified account")
	ErrLinkExpired   = errors.New("link expired")
	ErrInvalidTOTP   = errors.New("invalid TOTP")
)

var tmpl *template.Template

func SignUp(body *transfer.SignUpBody) (model.Account, error) {
	if _, err := model.FindAccount(body.Email); err != nil &&
		!errors.Is(err, mongo.ErrNoDocuments) {
		return model.Account{}, err
	} else if err == nil {
		return model.Account{}, ErrAccountExists
	}
	return transactionWithValue(
		func(ctx mongo.SessionContext) (model.Account, error) {
			account := model.Account{
				Email:    body.Email,
				Name:     body.Name,
				Password: argon2.Hash(body.Password),
				Role:     account.Member,
				Status:   account.Unverified,
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

func Verify(id string, a *model.Account) (model.Account, error) {
	if err := checkStatus(a); err != nil {
		return model.Account{}, err
	}
	if err := model.FindVerificationLink(id); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Account{}, ErrLinkExpired
		}
		return model.Account{}, err
	}
	return transactionWithValue(
		func(ctx mongo.SessionContext) (model.Account, error) {
			if err := model.DeleteVerificationLinks(ctx, a.Email); err != nil {
				return model.Account{}, err
			}
			return model.UpdateAccount(
				ctx,
				bson.M{"email": a.Email},
				bson.M{"$set": bson.M{"status": account.Normal}},
			)
		},
	)
}

func SendVerificationLink(account *model.Account) error {
	if err := checkStatus(account); err != nil {
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
		return model.Account{}, err
	}
	if !argon2.Verify(body.Password, account.Password) {
		return model.Account{}, ErrPassword
	}
	return account, nil
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
			image, err := key.Image(512, 512)
			if err != nil {
				return transfer.CreateTOTPBody{}, err
			}
			buf := bytes.Buffer{}
			if err := jpeg.Encode(&buf, image, nil); err != nil {
				return transfer.CreateTOTPBody{}, err
			}
			fmt.Println(key.String(), key.URL())
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
		return err
	}
	if !totp.Validate(code, t.Secret) {
		return ErrInvalidTOTP
	}
	return transaction(func(ctx mongo.SessionContext) error {
		if _, err := model.UpdateAccount(
			ctx,
			bson.M{"email": email},
			bson.M{"$set": bson.M{"totp": t.Secret}},
		); err != nil {
			return err
		}
		return model.DeleteTOTP(ctx, email)
	})
}

func DeleteTOTP(email string) error {
	return transaction(func(ctx mongo.SessionContext) error {
		if _, err := model.UpdateAccount(
			ctx,
			bson.M{"email": email},
			bson.M{"$unset": bson.M{"totp": ""}},
		); err != nil {
			return err
		}
		return model.DeleteTOTP(ctx, email)
	})
}

func generateCode(length int) string {
	digits := make([]byte, length)
	for i := range length {
		digits[i] = byte(rand.Intn(10)) + '0'
	}
	return string(digits)
}

func checkStatus(a *model.Account) error {
	if a.Status != account.Unverified {
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
