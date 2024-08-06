package service

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/bionicosmos/aegle/edge"
	pb "github.com/bionicosmos/aegle/edge/xray"
	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/service/subscription"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrUnsupportedAction = errors.New("unsupported action")

func FindUsers(page int) (Pagination[model.User], error) {
	return paginate(page, "users", model.FindUsers)
}

func UpdateUserProfiles(body *transfer.UpdateUserProfilesBody) error {
	actionMap := map[string]string{
		"add":    "$push",
		"remove": "$pull",
	}
	arrayAction, ok := actionMap[body.Action]
	if !ok {
		return ErrUnsupportedAction
	}
	return transaction(func(ctx mongo.SessionContext) error {
		profile, err := model.UpdateProfile(
			ctx,
			bson.M{"name": body.ProfileName},
			bson.M{arrayAction: bson.M{"userIds": body.Id}},
		)
		if err != nil {
			return err
		}
		user, err := model.FindUser(ctx, &body.Id)
		if err != nil {
			return err
		}
		arrayFilter := bson.M{}
		if body.Action == "add" {
			link, err := subscription.Generate(&profile, &user)
			if err != nil {
				return err
			}
			arrayFilter = bson.M{"name": body.ProfileName, "link": link}
		}
		if body.Action == "remove" {
			arrayFilter = bson.M{"name": body.ProfileName}
		}
		if err := model.UpdateUser(
			ctx,
			bson.M{"_id": body.Id},
			bson.M{arrayAction: bson.M{"profiles": arrayFilter}},
		); err != nil {
			return err
		}
		node, err := model.FindNode(ctx, &profile.NodeId)
		if err != nil {
			return err
		}
		if body.Action == "add" {
			return edge.AddUser(
				node.APIAddress,
				&pb.AddUserRequest{
					ProfileName: body.ProfileName,
					User: &pb.User{
						Email: user.Email,
						Level: user.Level,
						Uuid:  user.UUID,
						Flow:  user.Flow,
					},
				},
			)
		}
		return edge.RemoveUser(
			node.APIAddress,
			&pb.RemoveUserRequest{
				ProfileName: body.ProfileName,
				Email:       user.Email,
			},
		)
	})
}

func DeleteUser(id *primitive.ObjectID) error {
	return transaction(func(ctx mongo.SessionContext) error {
		user, err := model.DeleteUser(ctx, id)
		if err != nil {
			return err
		}
		wg := sync.WaitGroup{}
		errCh := make(chan error)
		for _, profile := range user.Profiles {
			profile, err := model.UpdateProfile(
				ctx,
				bson.M{"name": profile.Name},
				bson.M{"$pull": bson.M{"userIds": user.Id}},
			)
			if err != nil {
				return err
			}
			node, err := model.FindNode(ctx, &profile.NodeId)
			if err != nil {
				return err
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				errCh <- edge.RemoveUser(
					node.APIAddress,
					&pb.RemoveUserRequest{
						ProfileName: profile.Name,
						Email:       user.Email,
					},
				)
			}()
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
	})
}

func CheckUser() error {
	users, err := model.FindUsers(options.Find())
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	errCh := make(chan error)
	for _, user := range users {
		nextDate, err := time.Parse(
			time.RFC3339,
			strings.Split(user.NextDate, "[")[0],
		)
		if err != nil {
			return err
		}
		if nextDate.Before(time.Now()) {
			for _, profile := range user.Profiles {
				wg.Add(1)
				go func() {
					defer wg.Done()
					errCh <- UpdateUserProfiles(
						&transfer.UpdateUserProfilesBody{
							Id:          user.Id,
							ProfileName: profile.Name,
							Action:      "remove",
						},
					)
				}()
			}
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
