package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id        primitive.ObjectID    `json:"id,omitempty" bson:"_id"`
	Name      string                `json:"name,omitempty"`
	Email     string                `json:"email,omitempty"`
	Level     uint32                `json:"level,omitempty"`
	StartDate string                `json:"startDate,omitempty" bson:"startDate"`
	Cycles    int                   `json:"cycles,omitempty"`
	NextDate  string                `json:"nextDate,omitempty" bson:"nextDate"`
	UUID      string                `json:"uuid,omitempty" bson:"uuid"`
	Flow      string                `json:"flow,omitempty"`
	Security  string                `json:"security,omitempty"`
	Profiles  []ProfileSubscription `json:"profiles,omitempty"`
}

type ProfileSubscription struct {
	Name string `json:"name,omitempty"`
	Link string `json:"link,omitempty"`
}

var usersColl *mongo.Collection

func FindUser(id *primitive.ObjectID) (User, error) {
	user := User{}
	return user, usersColl.
		FindOne(context.Background(), bson.M{"_id": *id}).
		Decode(&user)
}

func FindUsers(query *Query) ([]User, error) {
	ctx := context.Background()
	cursor, err := usersColl.Find(
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
	var users []User
	return users, cursor.All(ctx, &users)
}

func InsertUser(ctx context.Context, user *User) error {
	user.Id = primitive.NewObjectID()
	_, err := usersColl.InsertOne(ctx, user)
	return err
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
