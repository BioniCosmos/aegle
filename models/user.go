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
	Id        primitive.ObjectID            `json:"id" bson:"_id"`
	Name      string                        `json:"name"`
	Email     string                        `json:"email"`
	Level     uint32                        `json:"level"`
	StartDate string                        `json:"startDate" bson:"startDate"`
	Cycles    int                           `json:"cycles" bson:"cycles"`
	NextDate  string                        `json:"nextDate" bson:"nextDate"`
	Account   map[string]json.RawMessage    `json:"account"`
	Profiles  map[primitive.ObjectID]string `json:"profiles"`
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

type UserResponse struct {
	User     User      `json:"user,omitempty"`
	Profiles []Profile `json:"profiles,omitempty"`
}

func UserProfilesAggregateQuery(skip int64, limit int64) ([]UserResponse, error) {
	pipeline := bson.A{
		bson.M{"$skip": skip},
		bson.M{"$limit": limit},
		bson.M{
			"$project": bson.M{
				"_id": 0,
				"user": bson.M{
					"_id":       "$_id",
					"account":   "$account",
					"startDate": "$startDate",
					"cycles":    "$cycles",
					"nextDate":  "$nextDate",
					"email":     "$email",
					"level":     "$level",
					"name":      "$name",
					"profiles":  "$profiles",
				},
				"profiles": bson.M{
					"$map": bson.M{
						"input": bson.M{"$objectToArray": "$profiles"},
						"as":    "profile",
						"in": bson.M{
							"$convert": bson.M{
								"input": "$$profile.k",
								"to":    "objectId",
							},
						},
					},
				},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "profiles",
				"localField":   "profiles",
				"foreignField": "_id",
				"as":           "profiles",
			},
		},
	}

	cursor, err := usersColl.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	res := make([]UserResponse, 0)
	if err := cursor.All(context.Background(), &res); err != nil {
		return nil, err
	}

	return res, nil
}

func QueryBillExpiredUsers() ([]User, error) {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"$expr": bson.M{
					"$lt": bson.A{
						bson.M{
							"$toDate": bson.M{
								"$first": bson.M{
									"$split": bson.A{"$nextDate", "["},
								},
							},
						},
						time.Now(),
					},
				},
			},
		},
	}
	cursor, err := usersColl.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)
	return users, cursor.All(context.TODO(), &users)
}
