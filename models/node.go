package models

import (
    "context"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Node struct {
    Id         primitive.ObjectID `json:"id" bson:"_id"`
    Name       string             `json:"name"`
    APIAddress string             `json:"apiAddress" bson:"apiAddress"`
    APIPort    int                `json:"apiPort" bson:"apiPort"`
}

var nodesColl *mongo.Collection

func FindNode(id string) (Node, error) {
    _id, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return Node{}, err
    }

    var node Node
    err = nodesColl.FindOne(context.TODO(), bson.D{
        {Key: "_id", Value: _id},
    }).Decode(&node)
    return node, err
}

func FindNodes(skip int64, limit int64) ([]Node, error) {
    cursor, err := nodesColl.Find(context.TODO(), bson.D{}, options.Find().SetSort(bson.D{
        {Key: "name", Value: 1},
    }).SetSkip(skip).SetLimit(limit))
    if err != nil {
        return nil, err
    }

    var nodes []Node
    err = cursor.All(context.TODO(), &nodes)
    return nodes, err
}

func (node *Node) Insert() error {
    node.Id = primitive.NewObjectID()
    _, err := nodesColl.InsertOne(context.TODO(), node)
    return err
}

func (node *Node) Update(id string) error {
    _id, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = nodesColl.UpdateOne(
        context.TODO(),
        bson.D{
            {Key: "_id", Value: _id},
        },
        bson.D{
            {Key: "$set", Value: node},
        },
    )
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

// func FindInboundsByTag(inboundTags ...string) ([]Inbound, error) {
//     cursor, err := nodesColl.Aggregate(context.TODO(), mongo.Pipeline{
//         {
//             {Key: "$unwind", Value: "$inbounds"},
//         },
//         {
//             {Key: "$match", Value: bson.D{
//                 {Key: "inbounds.tag", Value: bson.D{
//                     {Key: "$in", Value: inboundTags},
//                 }},
//             }},
//         },
//         {
//             {Key: "$project", Value: bson.D{
//                 {Key: "inbounds", Value: true},
//                 {Key: "_id", Value: false},
//             }},
//         },
//     })
//     if err != nil {
//         return nil, err
//     }

//     var inbounds []Inbound
//     err = cursor.All(context.TODO(), &inbounds)
//     if err != nil {
//         return nil, err
//     }
//     return inbounds, nil
// }
