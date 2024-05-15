package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Node struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name         string             `json:"name,omitempty"`
	APIAddress   string             `json:"apiAddress,omitempty" bson:"apiAddress"`
	ProfileNames []string           `json:"profileNames,omitempty" bson:"profileNames"`
}

var nodesColl *mongo.Collection

func FindNode(id *primitive.ObjectID) (Node, error) {
	node := Node{}
	return node, nodesColl.
		FindOne(context.Background(), bson.M{"_id": *id}).
		Decode(&node)
}

func FindNodes(query *Query) ([]Node, error) {
	ctx := context.Background()
	cursor, err := nodesColl.Find(
		ctx,
		bson.M{},
		options.
			Find().
			SetSort(bson.M{"name": 1}).
			SetSkip(query.Skip).
			SetLimit(query.Limit),
	)
	if err != nil {
		return nil, err
	}
	var nodes []Node
	return nodes, cursor.All(ctx, &nodes)
}

func InsertNode(ctx context.Context, node *Node) error {
	node.Id = primitive.NewObjectID()
	_, err := nodesColl.InsertOne(ctx, node)
	return err
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
