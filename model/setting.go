package model

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Setting struct {
	BaseURL string `bson:"baseURL"`
	Email   *Email
}

type Email struct {
	Host     string
	Port     uint16
	Username string
	Password string
}

var settingsColl *mongo.Collection

const id = "config"

func LoadSettings() (setting Setting, err error) {
	ctx := context.Background()
	if err = settingsColl.
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&setting); errors.Is(err, mongo.ErrNoDocuments) {
		_, err = settingsColl.InsertOne(ctx, bson.M{"_id": id})
	}
	return
}
