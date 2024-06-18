package service

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/model/account"
	"github.com/bionicosmos/argon2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(body *transfer.SignUpBody) error {
	return client.UseSession(
		context.Background(),
		func(ctx mongo.SessionContext) error {
			_, err := ctx.WithTransaction(
				ctx,
				func(ctx mongo.SessionContext) (any, error) {
					code := generateCode(4)
					model.InsertAccount(
						ctx,
						&model.Account{
							Email:    body.Email,
							Name:     body.Name,
							Password: argon2.Hash(body.Password),
							Role:     account.Member,
							Status:   account.Unverified,
						},
					)
					model.InsertVerificationCode(ctx, &model.VerificationCode{
						Email:     body.Email,
						Code:      code,
						CreatedAt: time.Now(),
					})
					if err := sendMail(
						body.Email,
						"Verification",
						code,
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

func Verify(body *transfer.VerifyBody) (bool, error) {
	verificationCode, err := model.FindVerificationCode(body.Email)
	if err != nil {
		return false, err
	}
	return body.Code == verificationCode.Code, nil
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
	setting, err := model.Setting("email")
	if err != nil {
		return err
	}
	email := setting.Email
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
