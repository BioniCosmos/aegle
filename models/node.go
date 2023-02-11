package models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Node struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name"`
	APIAddress string             `json:"apiAddress" bson:"apiAddress"`
}

var nodesColl *mongo.Collection

func FindNode(id any) (*Node, error) {
	var _id primitive.ObjectID
	switch id := id.(type) {
	case string:
		var err error
		_id, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
	case primitive.ObjectID:
		_id = id
	default:
		return nil, errors.New("invalid type for id")
	}

	var node Node
	return &node, nodesColl.FindOne(context.TODO(), bson.D{
		{Key: "_id", Value: _id},
	}).Decode(&node)
}

func FindNodes(skip int64, limit int64) ([]Node, error) {
	cursor, err := nodesColl.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{
		{Key: "name", Value: 1},
	}).SetSkip(skip).SetLimit(limit))
	if err != nil {
		return nil, err
	}

	var nodes []Node
	return nodes, cursor.All(context.TODO(), &nodes)
}

func (node *Node) Insert() error {
	node.Id = primitive.NewObjectID()
	_, err := nodesColl.InsertOne(context.TODO(), node)
	return err
}

func (node *Node) Update() error {
	_, err := nodesColl.UpdateByID(context.TODO(), node.Id, bson.D{
		{Key: "$set", Value: node},
	})
	return err
}

func DeleteNode(id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = nodesColl.DeleteOne(context.TODO(), bson.D{
		{Key: "_id", Value: _id},
	})
	return err
}
