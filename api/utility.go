package api

import (
    "context"
    "errors"
    "time"

    "github.com/bionicosmos/submgr/models"
    "go.mongodb.org/mongo-driver/bson"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

func dialAPIServer(apiAddress string) (*grpc.ClientConn, context.Context, context.CancelFunc, error) {
    cred, err := credentials.NewClientTLSFromFile("./cert.pem", "")
    if err != nil {
        return nil, nil, nil, err
    }
    ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
    conn, err := grpc.DialContext(
        ctx,
        apiAddress,
        grpc.WithTransportCredentials(cred),
    )
    if err != nil {
        return nil, nil, cancelFunc, errors.New("failed to dial " + apiAddress)
    }
    return conn, ctx, cancelFunc, nil
}

func CheckAllInbounds() error {
    profiles, err := models.FindProfiles(bson.D{}, bson.D{}, 0, 0)
    if err != nil {
        return err
    }
    for _, profile := range profiles {
        node, err := models.FindNode(profile.NodeId)
        if err != nil {
            return err
        }
        if err := AddInbound(profile.Inbound.ToConf(), node.APIAddress); err != nil {
            return err
        }
    }
    return nil
}

func CheckAllUsers() error {
    users, err := models.FindUsers(bson.D{}, bson.D{}, 0, 0)
    if err != nil {
        return err
    }
    for _, user := range users {
        for profileId := range user.Profiles {
            profile, err := models.FindProfile(profileId.Hex())
            if err != nil {
                return err
            }
            node, err := models.FindNode(profile.NodeId)
            if err != nil {
                return err
            }
            if err := AddUser(profile.Inbound.ToConf(), &user, node.APIAddress); err != nil {
                return err
            }
        }
    }
    return nil
}
