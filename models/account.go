package models

import (
    "bytes"
    "context"
    "runtime"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/argon2"
)

type Account struct {
    Username string `json:"username"`
    Password string `json:"password"`
    IsAdmin  bool   `json:"isAdmin"`
}

type innerAccount struct {
    Username string
    Password []byte
    IsAdmin  bool `bson:"isAdmin"`
}

var accountsColl *mongo.Collection

func (account *Account) Insert() error {
    _, err := accountsColl.InsertOne(context.TODO(), newInnerAccount(account))
    return err
}

func FindAccount(username string) (*innerAccount, error) {
    inner := new(innerAccount)
    return inner, accountsColl.FindOne(context.TODO(), bson.D{
        {Key: "username", Value: username},
    }).Decode(inner)
}

func (account *Account) PasswordIsCorrect(inner *innerAccount) bool {
    return bytes.Equal(deriveKey(account.Password), inner.Password)
}

func (account *Account) Update() error {
    _, err := accountsColl.UpdateOne(context.TODO(), bson.D{
        {Key: "username", Value: account.Username},
    }, bson.D{
        {Key: "$set", Value: newInnerAccount(account)},
    })
    return err
}

func DeleteAccount(username string) error {
    _, err := accountsColl.DeleteOne(context.TODO(), bson.D{
        {Key: "username", Value: username},
    })
    return err
}

func newInnerAccount(account *Account) *innerAccount {
    return &innerAccount{
        Username: account.Username,
        Password: deriveKey(account.Password),
        IsAdmin:  account.IsAdmin,
    }
}

func deriveKey(password string) []byte {
    return argon2.IDKey(
        []byte(password),
        []byte("&y90w6b2&oUEy*%u9Tjz"),
        1,
        64*1024,
        uint8(runtime.NumCPU()),
        32,
    )
}
