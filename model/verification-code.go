package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VerificationCode struct {
	Email     string
	Code      string
	CreatedAt time.Time `bson:"createdAt"`
}

var verificationCodesColl *mongo.Collection

func FindVerificationCode(email string) (VerificationCode, error) {
	verificationCode := VerificationCode{}
	return verificationCode,
		verificationCodesColl.
			FindOne(
				context.Background(),
				bson.M{"email": email},
				options.FindOne().SetSort(bson.M{"createdAt": -1}),
			).
			Decode(&verificationCode)
}

func InsertVerificationCode(
	ctx context.Context,
	verificationCode *VerificationCode,
) error {
	_, err := verificationCodesColl.InsertOne(ctx, verificationCode)
	return err
}
