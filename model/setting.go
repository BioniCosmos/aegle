package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func FindSetting(field string) (Setting, error) {
	setting := Setting{}
	return setting, settingsColl.FindOne(
		context.Background(),
		bson.M{"_id": id},
		options.FindOne().SetProjection(bson.M{field: 1}),
	).Decode(&setting)
}

func SetSetting(field string, value any) error {
	_, err := settingsColl.UpdateByID(
		context.Background(),
		id,
		bson.M{"$set": bson.M{field: value}},
		options.Update().SetUpsert(true),
	)
	return err
}
