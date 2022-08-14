package models

import (
    "context"
    "log"

    "github.com/bionicosmos/submgr/config"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func Init() {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Conf.DatabaseURL))
    if err != nil {
        log.Fatal(err)
    }
    db = client.Database("submgr")
    nodesColl = db.Collection("nodes")
}
