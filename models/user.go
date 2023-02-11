package models

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id          primitive.ObjectID            `json:"id" bson:"_id"`
	Name        string                        `json:"name"`
	Email       string                        `json:"email"`
	Level       uint32                        `json:"level"`
	BillingDate *time.Time                    `json:"billingDate" bson:"billingDate"`
	Account     map[string]json.RawMessage    `json:"account"`
	Profiles    map[primitive.ObjectID]string `json:"profiles"`
}

var usersColl *mongo.Collection

func FindUser(id string) (*User, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user User
	return &user, usersColl.FindOne(context.TODO(), bson.D{
		{Key: "_id", Value: _id},
	}).Decode(&user)
}

func FindUsers(filter any, sort any, skip int64, limit int64) ([]User, error) {
	cursor, err := usersColl.Find(context.TODO(), filter, options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit))
	if err != nil {
		return nil, err
	}

	var users []User
	return users, cursor.All(context.TODO(), &users)
}

func (user *User) Insert() error {
	user.Id = primitive.NewObjectID()
	_, err := usersColl.InsertOne(context.TODO(), user)
	return err
}

func (user *User) Update() error {
	_, err := usersColl.UpdateByID(context.TODO(), user.Id, bson.D{
		{Key: "$set", Value: user},
	})
	return err
}

func (user *User) RemoveProfile(profileId string) error {
	_, err := usersColl.UpdateByID(context.TODO(), user.Id, bson.D{
		{Key: "$unset", Value: bson.D{
			{Key: "profiles." + profileId, Value: ""},
		}},
	})
	return err
}

func DeleteUser(id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = usersColl.DeleteOne(context.TODO(), bson.D{
		{Key: "_id", Value: _id},
	})
	return err
}
