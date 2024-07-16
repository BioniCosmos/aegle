package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Node struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name"`
	APIAddress   string             `json:"apiAddress" bson:"apiAddress"`
	ProfileNames []string           `json:"profileNames" bson:"profileNames"`
}

var nodesColl *mongo.Collection

func FindNode(ctx context.Context, id *primitive.ObjectID) (Node, error) {
	node := Node{}
	return node, nodesColl.FindOne(ctx, bson.M{"_id": *id}).Decode(&node)
}

func FindNodes(options *options.FindOptions) ([]Node, error) {
	ctx := context.Background()
	cursor, err := nodesColl.Find(
		ctx,
		bson.M{},
		options.SetSort(bson.M{"name": 1}),
	)
	if err != nil {
		return nil, err
	}
	var nodes []Node
	return nodes, cursor.All(ctx, &nodes)
}

func InsertNode(ctx context.Context, node *Node) (primitive.ObjectID, error) {
	node.Id = primitive.NewObjectID()
	node.ProfileNames = make([]string, 0)
	_, err := nodesColl.InsertOne(ctx, node)
	return node.Id, err
}

func UpdateNode(ctx context.Context, filter any, update any) (Node, error) {
	node := Node{}
	return node, nodesColl.FindOneAndUpdate(ctx, filter, update).Decode(&node)
}

func DeleteNode(ctx context.Context, id *primitive.ObjectID) (Node, error) {
	node := Node{}
	return node, nodesColl.FindOneAndDelete(
		ctx,
		bson.M{"_id": *id},
	).Decode(&node)
}
