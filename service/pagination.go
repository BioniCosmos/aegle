package service

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

const pageSize = 10

func paginate[T any](
	page int,
	collection string,
	f func(*options.FindOptions) ([]T, error),
) (Pagination[T], error) {
	total, err := client.
		Database(os.Getenv("DB_NAME")).
		Collection(collection).
		EstimatedDocumentCount(context.Background())
	if err != nil {
		return Pagination[T]{}, err
	}
	items, err := f(
		options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(pageSize),
	)
	if err != nil {
		return Pagination[T]{}, err
	}
	if items == nil {
		items = make([]T, 0)
	}
	return Pagination[T]{
		Items: items,
		Total: int(total),
	}, nil
}
