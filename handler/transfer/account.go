package transfer

import (
	"github.com/bionicosmos/aegle/model"
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
	Email  string            `json:"email"`
	Name   string            `json:"name"`
	Role   model.AccountRole `json:"role"`
	Status AccountStatus     `json:"status"`
}

type AccountStatus string

const (
	AccountSignedIn   AccountStatus = "signedIn"
	AccountUnverified AccountStatus = "unverified"
	AccountNeedMFA    AccountStatus = "needMFA"
)

func ToAccount(account *model.Account, status AccountStatus) Account {
	return Account{
		Email:  account.Email,
		Name:   account.Name,
		Role:   account.Role,
		Status: status,
	}
}

type MFABody struct {
	Type MFAType `json:"type"`
	Code string  `json:"code"`
}

type MFAType string

const (
	MFATypeTOTP         MFAType = "totp"
	MFATypeRecoveryCode MFAType = "recoveryCode"
)

type CreateTOTPBody struct {
	Secret string `json:"secret"`
	Image  string `json:"image"`
}

type ConfirmTOTPBody struct {
	Code string `json:"code"`
}
