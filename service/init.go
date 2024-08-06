package service

import (
	"context"
	"html/template"

	"github.com/bionicosmos/aegle/email"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func Init(c *mongo.Client) {
	client = c
	tmpl = template.Must(template.New("email").Parse(email.Template))
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
				func(ctx mongo.SessionContext) (any, error) { return f(ctx) },
			)
			if v != nil {
				value = v.(T)
			}
			return err
		},
	)
}
