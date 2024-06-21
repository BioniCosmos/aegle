package service

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/model/account"
	"github.com/bionicosmos/aegle/setting"
	"github.com/bionicosmos/argon2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrAccountExists = errors.New("account exists")
var ErrPassword = errors.New("password mismatch")

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

func Verify(id string) (model.Account, error) {
	return transactionWithValue(
		func(ctx mongo.SessionContext) (model.Account, error) {
			link, err := model.DeleteVerificationLink(ctx, id)
			if err != nil {
				return model.Account{}, err
			}
			return model.UpdateAccount(
				ctx,
				bson.M{"email": link.Email},
				bson.M{"$set": bson.M{"status": account.Normal}},
			)
		},
	)
}

func SendVerificationLink(email string) error {
	link, err := model.FindVerificationLink(email)
	if err != nil {
		return err
	}
	return sendLinkEmail(&link)
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

func generateCode(length int) string {
	digits := make([]byte, length)
	for i := range length {
		digits[i] = byte(rand.Intn(10)) + '0'
	}
	return string(digits)
}

func sendLinkEmail(link *model.VerificationLink) error {
	return sendMail(
		link.Email,
		"Sign-up Verification",
		setting.X.BaseURL+"/dashboard/verification/"+link.Id,
	)
}

func sendMail(to string, subject string, body string) error {
	email := setting.X.Email
	auth := smtp.PlainAuth("", email.Username, email.Password, email.Host)
	message := bytes.Buffer{}
	message.WriteString("Subject: ")
	message.WriteString(subject)
	message.WriteString("\r\n\r\n")
	message.WriteString(body)
	message.WriteString("\r\n")
	return smtp.SendMail(
		fmt.Sprintf("%v:%v", email.Host, email.Port),
		auth,
		email.Username,
		[]string{to},
		message.Bytes(),
	)
}
