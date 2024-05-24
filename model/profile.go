package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Profile struct {
	Name     string               `json:"name,omitempty"`
	Outbound string               `json:"outbound,omitempty"`
	NodeId   primitive.ObjectID   `json:"nodeId,omitempty" bson:"nodeId"`
	UserIds  []primitive.ObjectID `json:"userIds,omitempty" bson:"userIds"`
}

var profilesColl *mongo.Collection

func FindProfile(id *primitive.ObjectID) (Profile, error) {
	profile := Profile{}
	return profile, profilesColl.
		FindOne(context.Background(), bson.M{"_id": *id}).
		Decode(&profile)
}

func FindProfiles() ([]Profile, error) {
	ctx := context.Background()
	cursor, err := profilesColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var profiles []Profile
	return profiles, cursor.All(ctx, &profiles)
}

func InsertProfile(ctx context.Context, profile *Profile) error {
	profile.UserIds = make([]primitive.ObjectID, 0)
	_, err := profilesColl.InsertOne(ctx, profile)
	return err
}

func UpdateProfile(
	ctx context.Context,
	filter any,
	update any,
) (Profile, error) {
	profile := Profile{}
	return profile, profilesColl.
		FindOneAndUpdate(ctx, filter, update).
		Decode(&profile)
}

func DeleteProfile(ctx context.Context, name string) (Profile, error) {
	profile := Profile{}
	return profile, profilesColl.
		FindOneAndDelete(ctx, bson.M{"name": name}).
		Decode(&profile)
}

func DeleteProfiles(ctx context.Context, filter any) ([]Profile, error) {
	cursor, err := profilesColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var profiles []Profile
	if err := cursor.All(ctx, &profiles); err != nil {
		return nil, err
	}
	_, err = profilesColl.DeleteMany(ctx, filter)
	return profiles, err
}
