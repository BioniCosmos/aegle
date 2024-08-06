package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id        primitive.ObjectID    `json:"id" bson:"_id"`
	Name      string                `json:"name"`
	Email     string                `json:"email"`
	Level     uint32                `json:"level"`
	StartDate time.Time             `json:"startDate" bson:"startDate"`
	Cycles    int                   `json:"cycles"`
	UUID      string                `json:"uuid" bson:"uuid"`
	Flow      string                `json:"flow"`
	Security  string                `json:"security"`
	Profiles  []ProfileSubscription `json:"profiles"`
}

type ProfileSubscription struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

var usersColl *mongo.Collection

func FindUser(ctx context.Context, id *primitive.ObjectID) (User, error) {
	user := User{}
	return user, usersColl.FindOne(ctx, bson.M{"_id": *id}).Decode(&user)
}

func FindUsers(options *options.FindOptions) ([]User, error) {
	ctx := context.Background()
	cursor, err := usersColl.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, err
	}
	users := []User{}
	return users, cursor.All(ctx, &users)
}

func InsertUser(ctx context.Context, user *User) (primitive.ObjectID, error) {
	user.Id = primitive.NewObjectID()
	user.Profiles = make([]ProfileSubscription, 0)
	_, err := usersColl.InsertOne(ctx, user)
	return user.Id, err
}

func UpdateUser(ctx context.Context, filter any, update any) error {
	_, err := usersColl.UpdateOne(ctx, filter, update)
	return err
}

func UpdateUsers(ctx context.Context, filter any, update any) error {
	_, err := usersColl.UpdateMany(ctx, filter, update)
	return err
}

func DeleteUser(ctx context.Context, id *primitive.ObjectID) (User, error) {
	user := User{}
	return user, usersColl.
		FindOneAndDelete(ctx, bson.M{"_id": *id}).
		Decode(&user)
}
