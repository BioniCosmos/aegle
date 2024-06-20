package service

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func Init(c *mongo.Client) {
	client = c
}

func transaction(f func(ctx mongo.SessionContext) error) error {
	_, err := transactionWithValue(func(ctx mongo.SessionContext) (any, error) {
		return nil, f(ctx)
	})
	return err
}

func transactionWithValue[T any](
	f func(ctx mongo.SessionContext) (T, error),
) (value T, err error) {
	return value, client.UseSession(
		context.Background(),
		func(ctx mongo.SessionContext) error {
			v, err := ctx.WithTransaction(
				ctx,
				func(ctx mongo.SessionContext) (any, error) {
					return f(ctx)
				},
			)
			if v != nil {
				value = v.(T)
			}
			return err
		},
	)
}
