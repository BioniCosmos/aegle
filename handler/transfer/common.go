package transfer

import "go.mongodb.org/mongo-driver/bson/primitive"

type ParamsId struct {
	Id primitive.ObjectID `params:"id"`
}
