package models

import (
    "context"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Inbound struct {
    Id     primitive.ObjectID `json:"id" bson:"_id"`
    Name   string             `json:"name"`
    Data   string             `json:"data"`
    NodeId primitive.ObjectID `json:"nodeId" bson:"nodeId"`
}

var inboundsColl *mongo.Collection

func FindInbound(id string) (*Inbound, error) {
    _id, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var inbound Inbound
    return &inbound, inboundsColl.FindOne(context.TODO(), bson.D{
        {Key: "_id", Value: _id},
    }).Decode(&inbound)
}

func FindInbounds(skip int64, limit int64) ([]Inbound, error) {
    cursor, err := inboundsColl.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{
        {Key: "name", Value: 1},
    }).SetSkip(skip).SetLimit(limit))
    if err != nil {
        return nil, err
    }

    var inbounds []Inbound
    return inbounds, cursor.All(context.TODO(), &inbounds)
}

func (inbound *Inbound) Insert() error {
    _, err := inboundsColl.InsertOne(context.TODO(), inbound)
    return err
}

func DeleteInbound(id string) error {
    _id, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = inboundsColl.DeleteOne(context.TODO(), bson.D{
        {Key: "_id", Value: _id},
    })
    return err
}
