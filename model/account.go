package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account struct {
	Email    string        `json:"email"`
	Name     string        `json:"name"`
	Password string        `json:"password"`
	Role     AccountRole   `json:"role"`
	Status   AccountStatus `json:"status"`
	TOTP     string        `json:"totp"`
}

type AccountRole string

const AccountRoleMember, AccountRoleAdmin AccountRole = "member", "admin"

type AccountStatus string

const (
	AccountStatusNormal     AccountStatus = "normal"
	AccountStatusUnverified AccountStatus = "unverified"
)

var accountsColl *mongo.Collection

func FindAccount(email string) (Account, error) {
	account := Account{}
	return account,
		accountsColl.
			FindOne(context.Background(), bson.M{"email": email}).
			Decode(&account)
}

func InsertAccount(ctx context.Context, account *Account) error {
	_, err := accountsColl.InsertOne(ctx, account)
	return err
}

func UpdateAccount(ctx context.Context, filter any, update any) error {
	_, err := accountsColl.UpdateOne(ctx, filter, update)
	return err
}
