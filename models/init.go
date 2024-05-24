package models

import (
	"context"
	"log"

	"github.com/bionicosmos/aegle/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init() *mongo.Client {
	ctx := context.Background()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(config.C.DatabaseURL),
	)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(config.C.DatabaseName)
	nodesColl = db.Collection("nodes")
	profilesColl = db.Collection("profiles")
	usersColl = db.Collection("users")
	accountsColl = db.Collection("accounts")

	options := options.Index().SetUnique(true)
	profilesColl.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{Keys: bson.M{"name": 1}, Options: options},
	)
	usersColl.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options},
	)
	accountsColl.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{Keys: bson.M{"username": 1}, Options: options},
	)
	return client
}
