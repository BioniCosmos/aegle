package transfer

import (
	"reflect"
	"strings"

	"github.com/bionicosmos/aegle/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindUserBodyFrom(user *model.User) map[string]any {
	value := reflect.ValueOf(*user)
	body := make(map[string]any)
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldName := typ.Field(i).Name
		if fieldName == "NextDate" {
			continue
		}
		if fieldName == "Profiles" {
			profileNames := make([]string, 0)
			for _, profile := range user.Profiles {
				profileNames = append(profileNames, profile.Name)
			}
			body["profileNames"] = profileNames
			continue
		}
		body[strings.Split(typ.Field(i).Tag.Get("json"), ",")[0]] = value.Field(i).Interface()
	}
	return body
}

type FindUserProfilesQuery struct {
	Id     primitive.ObjectID `query:"id"`
	Base64 bool               `query:"base64"`
}

type InsertUserBody struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Level     uint32 `json:"level,omitempty"`
	StartDate string `json:"startDate,omitempty"`
	Cycles    int    `json:"cycles,omitempty"`
	NextDate  string `json:"nextDate,omitempty"`
	UUID      string `json:"uuid,omitempty"`
	Flow      string `json:"flow,omitempty"`
	Security  string `json:"security,omitempty"`
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
	Id       primitive.ObjectID `json:"id,omitempty"`
	Cycles   int                `json:"cycles,omitempty"`
	NextDate string             `json:"nextDate,omitempty"`
}

type UpdateUserProfilesBody struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	ProfileName string             `json:"profileName,omitempty"`
	Action      string             `json:"action,omitempty"`
}
