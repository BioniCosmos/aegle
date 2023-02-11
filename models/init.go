package models

import (
    "context"
    "log"

    "github.com/bionicosmos/submgr/config"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Init() {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Conf.DatabaseURL))
    if err != nil {
        log.Fatal(err)
    }
    db = client.Database(config.Conf.DatabaseName)
    nodesColl = db.Collection("nodes")
    profilesColl = db.Collection("profiles")
    usersColl = db.Collection("users")
    accountsColl = db.Collection("accounts")

    accountsColl.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
        Keys: bson.D{
            {Key: "username", Value: 1},
        },
        Options: options.Index().SetUnique(true),
    })
}
