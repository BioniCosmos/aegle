package service

import (
	"github.com/bionicosmos/aegle/edge"
	pb "github.com/bionicosmos/aegle/edge/xray"
	"github.com/bionicosmos/aegle/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertProfile(profile *model.Profile, inbound string) error {
	return transaction(func(ctx mongo.SessionContext) error {
		if err := model.InsertProfile(ctx, profile); err != nil {
			return err
		}
		node, err := model.UpdateNode(
			ctx,
			bson.M{"_id": profile.NodeId},
			bson.M{"$push": bson.M{"profileNames": profile.Name}},
		)
		if err != nil {
			return err
		}
		return edge.AddInbound(
			node.APIAddress,
			&pb.AddInboundRequest{Name: profile.Name, Inbound: inbound},
		)
	})
}

func DeleteProfile(name string) error {
	return transaction(func(ctx mongo.SessionContext) error {
		profile, err := model.DeleteProfile(ctx, name)
		if err != nil {
			return err
		}
		node, err := model.UpdateNode(
			ctx,
			bson.M{"_id": profile.NodeId},
			bson.M{"$pull": bson.M{"profileNames": name}},
		)
		if err != nil {
			return err
		}
		err = model.UpdateUsers(
			ctx,
			bson.M{"_id": bson.M{"$in": profile.UserIds}},
			bson.M{"$pull": bson.M{"profiles": bson.M{"name": name}}},
		)
		if err != nil {
			return err
		}
		return edge.RemoveInbound(
			node.APIAddress, &pb.RemoveInboundRequest{Name: name},
		)
	})
}
