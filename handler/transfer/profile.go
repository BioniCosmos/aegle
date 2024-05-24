package transfer

import (
	"github.com/bionicosmos/aegle/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InsertProfileBody struct {
	Name     string             `json:"name,omitempty"`
	Inbound  string             `json:"inbound,omitempty"`
	Outbound string             `json:"outbound,omitempty"`
	NodeId   primitive.ObjectID `json:"nodeId,omitempty"`
}

func (body *InsertProfileBody) ToProfile() (*model.Profile, string) {
	return &model.Profile{
		Name:     body.Name,
		Outbound: body.Outbound,
		NodeId:   body.NodeId,
	}, body.Inbound
}
