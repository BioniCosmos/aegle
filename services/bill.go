package services

import (
	"sync"

	"github.com/bionicosmos/submgr/api"
	"github.com/bionicosmos/submgr/models"
	"github.com/bionicosmos/submgr/services/subscription"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckUserBill() error {
	users, err := models.QueryBillExpiredUsers()
	if err != nil {
		return err
	}

	chSize := 0
	for _, user := range users {
		chSize += len(user.Profiles)
	}
	wg := sync.WaitGroup{}
	errCh := make(chan error, chSize)

	for _, user := range users {
		profileIds := make([]primitive.ObjectID, 0, len(user.Profiles))
		for profileId := range user.Profiles {
			profileIds = append(profileIds, profileId)
		}
		profiles, err := models.FindProfiles(bson.M{
			"_id": bson.M{"$in": profileIds},
		}, bson.D{}, 0, 0)
		if err != nil {
			return err
		}

		for _, profile := range profiles {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := UserRemoveProfile(&user, &profile); err != nil {
					errCh <- err
				}
			}()
		}
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
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
	user.Profiles[profile.Id], err = subscription.Generate(profile, user.Account)
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
