package service

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"path"
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

func SignUp(body *transfer.SignUpBody) error {
	return client.UseSession(
		context.Background(),
		func(ctx mongo.SessionContext) error {
			_, err := ctx.WithTransaction(
				ctx,
				func(ctx mongo.SessionContext) (any, error) {
					if err := model.InsertAccount(
						ctx,
						&model.Account{
							Email:    body.Email,
							Name:     body.Name,
							Password: argon2.Hash(body.Password),
							Role:     account.Member,
							Status:   account.Unverified,
						},
					); err != nil {
						return nil, err
					}
					id := uuid.NewString()
					if err := model.InsertVerificationLink(
						ctx,
						&model.VerificationLink{
							Id:        id,
							Email:     body.Email,
							CreatedAt: time.Now(),
						},
					); err != nil {
						return nil, err
					}
					if err := sendMail(
						body.Email,
						"Sign-up Verification",
						path.Join(
							setting.Get[string]("BaseURL"),
							"dashboard/verification",
							id,
						),
					); err != nil {
						return nil, err
					}
					return nil, nil
				},
			)
			return err
		},
	)
}

func Verify(id string) (model.Account, error) {
	link, err := model.FindVerificationLink(id)
	if err != nil {
		return model.Account{}, err
	}
	return model.UpdateAccount(
		bson.M{"email": link.Email},
		bson.M{"$set": bson.M{"status": account.Normal}},
	)
}

func SignIn(body *transfer.SignInBody) (bool, error) {
	account, err := model.FindAccount(body.Email)
	if err != nil {
		return false, err
	}
	return argon2.Verify(body.Password, account.Password), nil
}

func generateCode(length int) string {
	digits := make([]byte, length)
	for i := range length {
		digits[i] = byte(rand.Intn(10)) + '0'
	}
	return string(digits)
}

func sendMail(to string, subject string, body string) error {
	email := setting.Get[*model.Email]("Email")
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
