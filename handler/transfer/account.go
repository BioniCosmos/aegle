package transfer

import (
	"time"

	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/model/account"
)

type SignUpBody struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Account struct {
	Email  string         `json:"email"`
	Name   string         `json:"name"`
	Role   account.Role   `json:"role"`
	Status account.Status `json:"status"`
	Expire time.Time      `json:"expire"`
}

func ToAccount(account *model.Account, maxAge time.Duration) Account {
	return Account{
		Email:  account.Email,
		Name:   account.Name,
		Role:   account.Role,
		Status: account.Status,
		Expire: time.Now().Add(maxAge).UTC(),
	}
}
