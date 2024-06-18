package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type setting struct {
	Email struct {
		Host     string
		Port     uint16
		Username string
		Password string
	}
}

var settingsColl *mongo.Collection

const id = "config"

func Setting(field string) (setting, error) {
	s := setting{}
	return s, settingsColl.FindOne(
		context.Background(),
		bson.M{"_id": id},
		options.FindOne().SetProjection(bson.M{field: 1}),
	).Decode(&s)
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
