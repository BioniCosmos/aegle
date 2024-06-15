package transfer

import (
	"github.com/bionicosmos/aegle/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FindUserProfilesQuery struct {
	Id     primitive.ObjectID `query:"id"`
	Base64 bool               `query:"base64"`
}

type InsertUserBody struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Level     uint32 `json:"level"`
	StartDate string `json:"startDate"`
	Cycles    int    `json:"cycles"`
	NextDate  string `json:"nextDate"`
	UUID      string `json:"uuid"`
	Flow      string `json:"flow"`
	Security  string `json:"security"`
}

func (body *InsertUserBody) ToUser() model.User {
	return model.User{
		Name:      body.Name,
		Email:     body.Email,
		Level:     body.Level,
		StartDate: body.StartDate,
		Cycles:    body.Cycles,
		NextDate:  body.NextDate,
		UUID:      body.UUID,
		Flow:      body.Flow,
		Security:  body.Security,
	}
}

type UpdateUserDateBody struct {
	Id       primitive.ObjectID `json:"id"`
	Cycles   int                `json:"cycles"`
	NextDate string             `json:"nextDate"`
}

type UpdateUserProfilesBody struct {
	Id          primitive.ObjectID `json:"id"`
	ProfileName string             `json:"profileName"`
	Action      string             `json:"action"`
}
