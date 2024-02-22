package api

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/bionicosmos/submgr/config"
	"github.com/bionicosmos/submgr/models"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func dialAPIServer(apiAddress string) (*grpc.ClientConn, context.Context, context.CancelFunc, error) {
	cred, err := credentials.NewClientTLSFromFile(config.Conf.Cert, "")
	if err != nil {
		return nil, nil, nil, errors.New("failed to read the certificate file")
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

func ResetNode(node *models.Node) error {
	profiles, err := models.FindProfiles(bson.D{{Key: "nodeId", Value: node.Id}}, bson.D{}, 0, 0)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errorCh := make(chan error)

	for _, profile := range profiles {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := AddInbound(profile.Inbound.ToConf(), node.APIAddress); err != nil {
				errorCh <- err
				return
			}

			users, err := models.FindUsers(bson.M{"profiles." + profile.Id.Hex(): bson.M{"$exists": true}}, bson.D{}, 0, 0)
			if err != nil {
				errorCh <- err
				return
			}

			for _, user := range users {
				if err := AddUser(profile.Inbound.ToConf(), &user, node.APIAddress); err != nil {
					errorCh <- err
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errorCh)
	}()

	for err := range errorCh {
		return err
	}
	return nil
}

func ResetAllNodes() error {
	nodes, err := models.FindNodes(0, 0)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error)

	for _, node := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ResetNode(&node)
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		log.Print(err)
	}
	return nil
}
