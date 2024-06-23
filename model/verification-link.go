package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type VerificationLink struct {
	Id        string `bson:"_id"`
	Email     string
	CreatedAt time.Time `bson:"createdAt"`
}

var verificationLinksColl *mongo.Collection

func FindVerificationLink(id string) error {
	return verificationLinksColl.
		FindOne(context.Background(), bson.M{"_id": id}).
		Err()
}

func InsertVerificationLink(
	ctx context.Context,
	verificationLink *VerificationLink,
) error {
	_, err := verificationLinksColl.InsertOne(ctx, verificationLink)
	return err
}

func DeleteVerificationLinks(ctx context.Context, email string) error {
	_, err := verificationLinksColl.DeleteMany(ctx, bson.M{"email": email})
	return err
}
