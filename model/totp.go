package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TOTP struct {
	Email     string
	Secret    string
	CreatedAt time.Time `bson:"createdAt"`
}

var totpColl *mongo.Collection

func FindTOTP(email string) (TOTP, error) {
	totp := TOTP{}
	return totp,
		totpColl.
			FindOne(context.Background(), bson.M{"email": email}).
			Decode(&totp)
}

func InsertTOTP(ctx context.Context, totp *TOTP) error {
	_, err := totpColl.InsertOne(ctx, totp)
	return err
}

func UpdateTOTP(ctx context.Context, filter, update any) error {
	_, err := totpColl.UpdateOne(ctx, filter, update)
	return err
}

func DeleteTOTP(ctx context.Context, email string) error {
	_, err := totpColl.DeleteOne(ctx, bson.M{"email": email})
	return err
}
