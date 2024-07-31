package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecoveryCode struct {
	Code  string `json:"code"`
	Used  bool   `json:"used"`
	Email string `json:"email"`
}

var recoveryCodesColl *mongo.Collection

func FindRecoveryCode(ctx context.Context, filter any) (RecoveryCode, error) {
	code := RecoveryCode{}
	return code, recoveryCodesColl.FindOne(ctx, filter).Decode(&code)
}

func FindRecoveryCodes(
	ctx context.Context,
	email string,
) ([]RecoveryCode, error) {
	cursor, err := recoveryCodesColl.Find(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	var codes []RecoveryCode
	return codes, cursor.All(ctx, &codes)
}

func CountRecoveryCodes(ctx context.Context, email string) (int64, error) {
	return recoveryCodesColl.CountDocuments(ctx, bson.M{"email": email})
}

func InsertRecoveryCodes(ctx context.Context, codes []RecoveryCode) error {
	docs := []any{}
	for _, code := range codes {
		docs = append(docs, code)
	}
	_, err := recoveryCodesColl.InsertMany(ctx, docs)
	return err
}

func UpdateRecoveryCode(ctx context.Context, filter any, update any) error {
	_, err := recoveryCodesColl.UpdateOne(ctx, filter, update)
	return err
}

func DeleteRecoveryCodes(ctx context.Context, email string) error {
	_, err := recoveryCodesColl.DeleteMany(ctx, bson.M{"email": email})
	return err
}
