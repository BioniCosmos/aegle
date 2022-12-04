package services

import (
    "time"

    "github.com/bionicosmos/submgr/api"
    "github.com/bionicosmos/submgr/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckUserBill() error {
    users, err := models.FindUsers(bson.D{
        {Key: "billingDate", Value: bson.D{
            {Key: "$lte", Value: time.Now().AddDate(0, -1, 0)},
        }},
    }, bson.D{}, 0, 0)
    if err != nil {
        return err
    }
    for _, user := range users {
        profileIds := make([]primitive.ObjectID, 0, len(user.Profiles))
        for profileId := range user.Profiles {
            profileIds = append(profileIds, profileId)
        }
        profiles, err := models.FindProfiles(bson.D{
            {Key: "_id", Value: bson.D{
                {Key: "$in", Value: profileIds},
            }},
        }, bson.D{}, 0, 0)
        if err != nil {
            return err
        }
        for _, profile := range profiles {
            if err := UserRemoveProfile(&user, &profile); err != nil {
                return err
            }
        }
    }
    return nil
}

func UserAddProfile(user *models.User, profile *models.Profile) error {
    node, err := models.FindNode(profile.NodeId)
    if err != nil {
        return err
    }
    if err := api.AddUser(profile.Inbound.ToConf(), user, node.APIAddress); err != nil {
        return err
    }
    user.Profiles[profile.Id], err = GenerateSubscription(user, profile)
    if err != nil {
        return err
    }
    return user.Update()
}

func UserRemoveProfile(user *models.User, profile *models.Profile) error {
    node, err := models.FindNode(profile.NodeId)
    if err != nil {
        return err
    }
    if err := api.RemoveUser(profile.Inbound.ToConf(), user, node.APIAddress); err != nil {
        return err
    }
    return user.RemoveProfile(profile.Id.Hex())
}
