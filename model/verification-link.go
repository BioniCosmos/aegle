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

func FindVerificationLink(id string) (VerificationLink, error) {
	verificationLink := VerificationLink{}
	return verificationLink,
		verificationLinksColl.
			FindOne(context.Background(), bson.M{"_id": id}).
			Decode(&verificationLink)
}

func InsertVerificationLink(
	ctx context.Context,
	verificationLink *VerificationLink,
) error {
	_, err := verificationLinksColl.InsertOne(ctx, verificationLink)
	return err
}
