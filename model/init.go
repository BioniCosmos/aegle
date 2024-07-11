package model

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init() *mongo.Client {
	ctx := context.Background()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(os.Getenv("DB_URL")),
	)
	if err != nil {
		panic(err)
	}
	db := client.Database(os.Getenv("DB_NAME"))
	nodesColl = db.Collection("nodes")
	profilesColl = db.Collection("profiles")
	usersColl = db.Collection("users")
	accountsColl = db.Collection("accounts")
	settingsColl = db.Collection("settings")
	verificationCodesColl = db.Collection("verificationCodes")
	verificationLinksColl = db.Collection("verificationLinks")
	totpColl = db.Collection("totp")

	profilesColl.Indexes().CreateOne(ctx, uniqueIndex("name"))
	usersColl.Indexes().CreateOne(ctx, uniqueIndex("email"))
	accountsColl.Indexes().CreateOne(ctx, uniqueIndex("email"))
	verificationCodesColl.
		Indexes().
		CreateOne(ctx, expireIndex("createdAt", 5*minute))
	verificationLinksColl.
		Indexes().
		CreateOne(ctx, expireIndex("createdAt", day))
	totpColl.Indexes().CreateOne(ctx, expireIndex("createdAt", 5*minute))
	return client
}

func uniqueIndex(field string) mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetUnique(true),
	}
}

const minute = 60
const hour = 60 * minute
const day = 24 * hour

func expireIndex(field string, seconds int32) mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetExpireAfterSeconds(seconds),
	}
}
