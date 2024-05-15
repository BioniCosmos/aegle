package edge

import (
	"context"

	pb "github.com/bionicosmos/aegle/edge/xray"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AddInbound(target string, in *pb.AddInboundRequest) error {
	client, err := newClient(target)
	if err != nil {
		return err
	}
	_, err = client.AddInbound(context.Background(), in)
	return err
}

func RemoveInbound(target string, in *pb.RemoveInboundRequest) error {
	client, err := newClient(target)
	if err != nil {
		return err
	}
	_, err = client.RemoveInbound(context.Background(), in)
	return err
}

func AddUser(target string, in *pb.AddUserRequest) error {
	client, err := newClient(target)
	if err != nil {
		return err
	}
	_, err = client.AddUser(context.Background(), in)
	return err
}

func RemoveUser(target string, in *pb.RemoveUserRequest) error {
	client, err := newClient(target)
	if err != nil {
		return err
	}
	_, err = client.RemoveUser(context.Background(), in)
	return err
}

func newClient(target string) (pb.XrayClient, error) {
	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return pb.NewXrayClient(conn), nil
}
