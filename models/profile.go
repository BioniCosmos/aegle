package models

import (
    "context"
    "encoding/json"

    "github.com/xtls/xray-core/common/net"
    "github.com/xtls/xray-core/infra/conf"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Profile struct {
    Id       primitive.ObjectID         `json:"id" bson:"_id"`
    Name     string                     `json:"name"`
    Inbounds []inboundDetourConfig      `json:"inbounds"`
    Outbound *conf.OutboundDetourConfig `json:"outbound"`
    NodeId   primitive.ObjectID         `json:"nodeId" bson:"nodeId"`
}

type inboundDetourConfig struct {
    Protocol       string                              `json:"protocol"`
    PortList       *conf.PortList                      `json:"port"`
    ListenOn       string                              `json:"listen"`
    Settings       *json.RawMessage                    `json:"settings"`
    Tag            string                              `json:"tag"`
    Allocation     *conf.InboundDetourAllocationConfig `json:"allocate"`
    StreamSetting  *conf.StreamConfig                  `json:"streamSettings"`
    DomainOverride *conf.StringList                    `json:"domainOverride"`
    SniffingConfig *conf.SniffingConfig                `json:"sniffing"`
}

var profilesColl *mongo.Collection

func FindProfile(id string) (*Profile, error) {
    _id, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var profile Profile
    return &profile, profilesColl.FindOne(context.TODO(), bson.D{
        {Key: "_id", Value: _id},
    }).Decode(&profile)
}

func FindProfiles(filter any, sort any, skip int64, limit int64) ([]Profile, error) {
    cursor, err := profilesColl.Find(context.TODO(), filter, options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit))
    if err != nil {
        return nil, err
    }

    var profiles []Profile
    return profiles, cursor.All(context.TODO(), &profiles)
}

func (profile *Profile) Insert() error {
    _, err := profilesColl.InsertOne(context.TODO(), profile)
    return err
}

func DeleteProfile(id string) error {
    _id, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }
    _, err = profilesColl.DeleteOne(context.TODO(), bson.D{
        {Key: "_id", Value: _id},
    })
    return err
}

func (inbound *inboundDetourConfig) ToConf() *conf.InboundDetourConfig {
    listenOn := new(conf.Address)
    if inbound.ListenOn == "" {
        inbound.ListenOn = "0.0.0.0"
    }
    listenOn.Address = net.ParseAddress(inbound.ListenOn)
    return &conf.InboundDetourConfig{
        Protocol:       inbound.Protocol,
        PortList:       inbound.PortList,
        ListenOn:       listenOn,
        Settings:       inbound.Settings,
        Tag:            inbound.Tag,
        Allocation:     inbound.Allocation,
        StreamSetting:  inbound.StreamSetting,
        DomainOverride: inbound.DomainOverride,
        SniffingConfig: inbound.SniffingConfig,
    }
}
