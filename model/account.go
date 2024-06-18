package model

import (
	"context"

	"github.com/bionicosmos/aegle/model/account"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Account struct {
	Email    string         `json:"email"`
	Name     string         `json:"name"`
	Password string         `json:"password"`
	Role     account.Role   `json:"role"`
	Status   account.Status `json:"status"`
}

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
