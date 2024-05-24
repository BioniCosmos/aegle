package service

import (
	"context"

	"github.com/bionicosmos/aegle/edge"
	pb "github.com/bionicosmos/aegle/edge/xray"
	"github.com/bionicosmos/aegle/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertProfile(profile *model.Profile, inbound string) error {
	ctx := context.Background()
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(
		ctx,
		func(ctx mongo.SessionContext) (interface{}, error) {
			if err := model.InsertProfile(ctx, profile); err != nil {
				return nil, err
			}
			node, err := model.UpdateNode(
				ctx,
				bson.M{"_id": profile.NodeId},
				bson.M{"$push": bson.M{"profileNames": profile.Name}},
			)
			if err != nil {
				return nil, err
			}
			return nil, edge.AddInbound(
				node.APIAddress,
				&pb.AddInboundRequest{
					Name:    profile.Name,
					Inbound: inbound,
				},
			)
		},
	)
	return err
}

func DeleteProfile(name string) error {
	ctx := context.Background()
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(
		ctx,
		func(ctx mongo.SessionContext) (interface{}, error) {
			profile, err := model.DeleteProfile(ctx, name)
			if err != nil {
				return nil, err
			}
			node, err := model.UpdateNode(
				ctx,
				bson.M{"_id": profile.NodeId},
				bson.M{"$pull": bson.M{"profileNames": name}},
			)
			if err != nil {
				return nil, err
			}
			err = model.UpdateUsers(
				ctx,
				bson.M{"_id": bson.M{"$in": profile.UserIds}},
				bson.M{"$pull": bson.M{"profiles": bson.M{"name": name}}},
			)
			if err != nil {
				return nil, err
			}
			return nil, edge.RemoveInbound(
				node.APIAddress, &pb.RemoveInboundRequest{Name: name},
			)
		},
	)
	return err
}
