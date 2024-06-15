package transfer

import (
	"github.com/bionicosmos/aegle/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertProfileBody struct {
	Name     string             `json:"name"`
	Inbound  string             `json:"inbound"`
	Outbound string             `json:"outbound"`
	NodeId   primitive.ObjectID `json:"nodeId"`
}

func (body *InsertProfileBody) ToProfile() (*model.Profile, string) {
	return &model.Profile{
		Name:     body.Name,
		Outbound: body.Outbound,
		NodeId:   body.NodeId,
	}, body.Inbound
}
