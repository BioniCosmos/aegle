package model

import (
	"context"

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

func FindSetting() (Setting, error) {
	setting := Setting{}
	return setting,
		settingsColl.
			FindOne(context.Background(), bson.M{"_id": id}).
			Decode(&setting)
}
