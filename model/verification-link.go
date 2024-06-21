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

func FindVerificationLink(email string) (VerificationLink, error) {
	link := VerificationLink{}
	return link,
		verificationLinksColl.
			FindOne(context.Background(), bson.M{"email": email}).
			Decode(&link)
}

func InsertVerificationLink(
	ctx context.Context,
	verificationLink *VerificationLink,
) error {
	_, err := verificationLinksColl.InsertOne(ctx, verificationLink)
	return err
}

func DeleteVerificationLink(
	ctx context.Context,
	id string,
) (VerificationLink, error) {
	link := VerificationLink{}
	return link,
		verificationLinksColl.
			FindOneAndDelete(ctx, bson.M{"_id": id}).
			Decode(&link)
}
