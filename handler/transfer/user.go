package transfer

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FindUserProfilesQuery struct {
	Id     primitive.ObjectID `query:"id"`
	Base64 bool               `query:"base64"`
}

type UpdateUserDateBody struct {
	Id     primitive.ObjectID `json:"id"`
	Cycles int                `json:"cycles"`
}

type UpdateUserProfilesBody struct {
	Id          primitive.ObjectID `json:"id"`
	ProfileName string             `json:"profileName"`
	Action      string             `json:"action"`
}
