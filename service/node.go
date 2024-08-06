package service

import (
	"sync"

	"github.com/bionicosmos/aegle/edge"
	pb "github.com/bionicosmos/aegle/edge/xray"
	"github.com/bionicosmos/aegle/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindNodes(page int) (Pagination[model.Node], error) {
	return paginate(page, "nodes", model.FindNodes)
}

func DeleteNode(id *primitive.ObjectID) error {
	return transaction(func(ctx mongo.SessionContext) error {
		node, err := model.DeleteNode(ctx, id)
		if err != nil {
			return err
		}
		profiles, err := model.DeleteProfiles(
			ctx,
			bson.M{"name": bson.M{"$in": node.ProfileNames}},
		)
		if err != nil {
			return err
		}
		for _, profile := range profiles {
			if err := model.UpdateUsers(
				ctx,
				bson.M{"_id": bson.M{"$in": profile.UserIds}},
				bson.M{
					"$pull": bson.M{"profiles": bson.M{"name": profile.Name}},
				},
			); err != nil {
				return err
			}
		}
		wg := sync.WaitGroup{}
		errCh := make(chan error)
		for _, profileName := range node.ProfileNames {
			wg.Add(1)
			go func() {
				defer wg.Done()
				errCh <- edge.RemoveInbound(
					node.APIAddress,
					&pb.RemoveInboundRequest{Name: profileName},
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
