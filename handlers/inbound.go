package handlers

import (
    "context"
    "encoding/json"
    "errors"
    "time"

    "github.com/bionicosmos/submgr/models"
    "github.com/gofiber/fiber/v2"
    "github.com/xtls/xray-core/app/proxyman/command"
    "github.com/xtls/xray-core/infra/conf"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func FindInbound(c *fiber.Ctx) error {
    inbound, err := models.FindInbound(c.Params("id"))
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.ErrNotFound
        }
        return err
    }
    return c.JSON(inbound)
}

func FindInbounds(c *fiber.Ctx) error {
    var query struct {
        skip  int64
        limit int64
    }
    if err := c.QueryParser(&query); err != nil {
        return err
    }

    inbounds, err := models.FindInbounds(query.skip, query.limit)
    if err != nil {
        return err
    }
    if inbounds == nil {
        return fiber.ErrNotFound
    }
    return c.JSON(inbounds)
}

func InsertInbound(c *fiber.Ctx) error {
    var inbound models.Inbound
    if err := c.BodyParser(&inbound); err != nil {
        return err
    }
    inbound.Id = primitive.NewObjectID()

    var inboundConfig conf.InboundDetourConfig
    if err := json.Unmarshal([]byte(inbound.Data), &inboundConfig); err != nil {
        return err
    }
    inboundConfig.Tag = inbound.Id.Hex()

    node, err := models.FindNode(inbound.NodeId)
    if err != nil {
        return err
    }

    conn, ctx, cancelFunc, err := dialAPIServer(node.APIAddress)
    defer cancelFunc()
    if err != nil {
        return err
    }

    client := command.NewHandlerServiceClient(conn)
    inboundHandlerConfig, err := inboundConfig.Build()
    if err != nil {
        return err
    }

    req := command.AddInboundRequest{
        Inbound: inboundHandlerConfig,
    }
    if _, err := client.AddInbound(ctx, &req); err != nil {
        return err
    }
    if err := conn.Close(); err != nil {
        return err
    }
    if err := inbound.Insert(); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusCreated)
}

func DeleteInbound(c *fiber.Ctx) error {
    inbound, err := models.FindInbound(c.Params("id"))
    if err != nil {
        return err
    }

    node, err := models.FindNode(inbound.NodeId)
    if err != nil {
        return err
    }
    if err := models.DeleteInbound(c.Params("id")); err != nil {
        return err
    }

    conn, ctx, cancelFunc, err := dialAPIServer(node.APIAddress)
    defer cancelFunc()
    if err != nil {
        return err
    }

    client := command.NewHandlerServiceClient(conn)
    req := command.RemoveInboundRequest{
        Tag: inbound.Id.Hex(),
    }
    if _, err := client.RemoveInbound(ctx, &req); err != nil {
        return err
    }
    if err := conn.Close(); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusNoContent)
}

func dialAPIServer(apiAddress string) (*grpc.ClientConn, context.Context, context.CancelFunc, error) {
    ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
    conn, err := grpc.DialContext(
        ctx,
        apiAddress,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
    )
    if err != nil {
        return nil, nil, cancelFunc, errors.New("failed to dial " + apiAddress)
    }
    return conn, ctx, cancelFunc, nil
}
