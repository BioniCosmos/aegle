package service

import "go.mongodb.org/mongo-driver/mongo"

var client *mongo.Client

func Init(c *mongo.Client) {
	client = c
}
