package service

import (
	"context"

	"github.com/bionicosmos/aegle/config"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination[T any] struct {
	Items        []T  `json:"items,omitempty"`
	CurrentPage  int  `json:"currentPage,omitempty"`
	PageSize     int  `json:"pageSize,omitempty"`
	TotalPages   int  `json:"totalPages,omitempty"`
	HasPrevious  bool `json:"hasPrevious,omitempty"`
	HasNext      bool `json:"hasNext,omitempty"`
	PreviousPage int  `json:"previousPage,omitempty"`
	NextPage     int  `json:"nextPage,omitempty"`
	FirstPage    int  `json:"firstPage,omitempty"`
	LastPage     int  `json:"lastPage,omitempty"`
}

const pageSize = 10

func paginate[T any](
	page int,
	collection string,
	f func(*options.FindOptions) ([]T, error),
) (Pagination[T], error) {
	total, err := client.
		Database(config.C.DatabaseName).
		Collection(collection).
		EstimatedDocumentCount(context.Background())
	if err != nil {
		return Pagination[T]{}, err
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	items, err := f(
		options.Find().SetSkip(int64((page - 1) * pageSize)).SetLimit(pageSize),
	)
	if err != nil {
		return Pagination[T]{}, err
	}
	return Pagination[T]{
		Items:        items,
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
		HasPrevious:  page > 1,
		HasNext:      page < totalPages,
		PreviousPage: max(page-1, 1),
		NextPage:     min(page+1, totalPages),
	}, nil
}
