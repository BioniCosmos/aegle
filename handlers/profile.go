package handlers

import (
    "fmt"

    "github.com/bionicosmos/submgr/api"
    "github.com/bionicosmos/submgr/models"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

func FindProfile(c *fiber.Ctx) error {
    profile, err := models.FindProfile(c.Params("id"))
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.ErrNotFound
        }
        return err
    }
    return c.JSON(profile)
}

func FindProfiles(c *fiber.Ctx) error {
    var query struct {
        Skip   int64              `query:"skip"`
        Limit  int64              `query:"limit"`
        NodeId primitive.ObjectID `query:"nodeId"`
    }
    if err := c.QueryParser(&query); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }

    var filter bson.E
    if !query.NodeId.IsZero() {
        filter = bson.E{Key: "nodeId", Value: query.NodeId}
    }
    profiles, err := models.FindProfiles(bson.D{
        filter,
    }, bson.D{
        {Key: "name", Value: 1},
    }, query.Skip, query.Limit)
    if err != nil {
        return err
    }
    if profiles == nil {
        return fiber.ErrNotFound
    }
    return c.JSON(profiles)
}

func InsertProfile(c *fiber.Ctx) error {
    var profile models.Profile
    if err := c.BodyParser(&profile); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }
    profile.Id = primitive.NewObjectID()

    // Find the node whose id equals to the nodeId.
    // If errors, return.
    // However, return 'node not found' if the node does not exist.
    node, err := models.FindNode(profile.NodeId)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.NewError(fiber.StatusBadRequest, "node not found")
        }
        return err
    }
    for i := range profile.Inbounds {
        profile.Inbounds[i].Tag = fmt.Sprintf("%v-%v", profile.Id.Hex(), i)
        if err := api.AddInbound(profile.Inbounds[i].ToConf(), node.APIAddress); err != nil {
            return err
        }
    }
    if err := profile.Insert(); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusCreated)
}

func DeleteProfile(c *fiber.Ctx) error {
    id := c.Params("id")
    if users, err := models.FindUsers(bson.D{
        {Key: fmt.Sprintf("profiles.%v", id), Value: bson.D{
            {Key: "$exists", Value: true},
        }},
    }, bson.D{}, 0, 0); err != nil {
        return err
    } else if users != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Users binding to the profile are not empty.")
    }

    profile, err := models.FindProfile(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.NewError(fiber.StatusBadRequest, "profile not found")
        }
        return err
    }

    node, err := models.FindNode(profile.NodeId)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.NewError(fiber.StatusBadRequest, "node not found")
        }
        return err
    }
    for _, inbound := range profile.Inbounds {
        if err := api.RemoveInbound(inbound.Tag, node.APIAddress); err != nil {
            return err
        }
    }
    if err := models.DeleteProfile(id); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusNoContent)
}
