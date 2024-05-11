package transfer

import (
	"encoding/json"

	"github.com/bionicosmos/aegle/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserGet struct {
	Id         primitive.ObjectID         `json:"id"`
	Name       string                     `json:"name"`
	Email      string                     `json:"email"`
	Level      uint32                     `json:"level"`
	StartDate  string                     `json:"startDate"`
	Cycles     int                        `json:"cycles"`
	Account    map[string]json.RawMessage `json:"account"`
	ProfileIds []primitive.ObjectID       `json:"profileIds"`
}

func UserGetFromUser(user *models.User) UserGet {
	profileIds := make([]primitive.ObjectID, 0)
	for profileId := range user.Profiles {
		profileIds = append(profileIds, profileId)
	}
	return UserGet{
		Id:         user.Id,
		Name:       user.Name,
		Email:      user.Email,
		Level:      user.Level,
		StartDate:  user.StartDate,
		Cycles:     user.Cycles,
		Account:    user.Account,
		ProfileIds: profileIds,
	}
}

type UserPost struct {
	Name       string                     `json:"name"`
	Email      string                     `json:"email"`
	Level      uint32                     `json:"level"`
	StartDate  string                     `json:"startDate"`
	Cycles     int                        `json:"cycles"`
	NextDate   string                     `json:"nextDate"`
	Account    map[string]json.RawMessage `json:"account"`
	ProfileIds []primitive.ObjectID       `json:"profileIds"`
}

func (userPost *UserPost) ToUser() models.User {
	return models.User{
		Name:      userPost.Name,
		Email:     userPost.Email,
		Level:     userPost.Level,
		StartDate: userPost.StartDate,
		Cycles:    userPost.Cycles,
		NextDate:  userPost.NextDate,
		Account:   userPost.Account,
		Profiles:  make(map[primitive.ObjectID]string),
	}
}

type UserPut struct {
	Id       string `json:"id"`
	Cycles   int    `json:"cycles"`
	NextDate string `json:"nextDate"`
}
