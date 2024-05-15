package transfer

import (
	"reflect"

	"github.com/bionicosmos/aegle/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindUserBodyFrom(user *models.User) map[string]any {
	value := reflect.ValueOf(user)
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
		body[fieldName] = value.Field(i).Interface()
	}
	return body
}

type FindUserProfilesQuery struct {
	Id     primitive.ObjectID `query:"id"`
	Base64 bool               `query:"base64"`
}

type InsertUserBody map[string]any

func (body *InsertUserBody) ToUser() models.User {
	user := models.User{}
	value := reflect.ValueOf(user)
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldName := typ.Field(i).Name
		if fieldName == "id" || fieldName == "Profiles" {
			continue
		}
		value.Field(i).Set(reflect.ValueOf((*body)[fieldName]))
	}
	return user
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
