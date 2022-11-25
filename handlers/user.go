package handlers

import (
    "encoding/base64"
    "strings"

    "github.com/bionicosmos/submgr/api"
    "github.com/bionicosmos/submgr/models"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

func FindUser(c *fiber.Ctx) error {
    user, err := models.FindUser(c.Params("id"))
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.ErrNotFound
        }
        return err
    }
    return c.JSON(user)
}

func FindUsers(c *fiber.Ctx) error {
    var query struct {
        Skip  int64 `query:"skip"`
        Limit int64 `query:"limit"`
    }
    if err := c.QueryParser(&query); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }

    users, err := models.FindUsers(bson.D{}, bson.D{
        {Key: "name", Value: 1},
    }, query.Skip, query.Limit)
    if err != nil {
        return err
    }
    if users == nil {
        return fiber.ErrNotFound
    }
    return c.JSON(users)
}

func InsertUser(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }
    for profileId := range user.Profiles {
        profile, err := models.FindProfile(profileId.Hex())
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
            if err := api.AddUser(inbound.ToConf(), &user, node.APIAddress); err != nil {
                return err
            }
        }
        user.Profiles[profileId], err = user.GenerateSubscription(profile)
        if err != nil {
            return fiber.NewError(fiber.StatusBadRequest, err.Error())
        }
    }
    if err := user.Insert(); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusCreated)
}

func UpdateUser(c *fiber.Ctx) error {
    id := c.Params("id")
    type operation uint
    const (
        Add uint = iota
        Remove
    )
    var body struct {
        Id        primitive.ObjectID `json:"id"`
        Operation operation          `json:"operation"`
    }
    if err := c.BodyParser(&body); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, err.Error())
    }
    user, err := models.FindUser(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.NewError(fiber.StatusBadRequest, "user not found")
        }
        return err
    }
    profile, err := models.FindProfile(body.Id.Hex())
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
        if body.Operation == operation(Add) {
            if err := api.AddUser(inbound.ToConf(), user, node.APIAddress); err != nil {
                return err
            }
            user.Profiles[body.Id], err = user.GenerateSubscription(profile)
            if err != nil {
                return fiber.NewError(fiber.StatusBadRequest, err.Error())
            }
            if err := user.Update(); err != nil {
                return err
            }
        } else {
            if err := api.RemoveUser(inbound.ToConf(), user, node.APIAddress); err != nil {
                return err
            }
            if err := user.RemoveProfile(body.Id.Hex()); err != nil {
                return err
            }
        }
    }
    return c.SendStatus(fiber.StatusNoContent)
}

func DeleteUser(c *fiber.Ctx) error {
    id := c.Params("id")
    user, err := models.FindUser(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.NewError(fiber.StatusBadRequest, "user not found")
        }
        return err
    }
    for profileId := range user.Profiles {
        profile, err := models.FindProfile(profileId.Hex())
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
            if err := api.RemoveUser(inbound.ToConf(), user, node.APIAddress); err != nil {
                return err
            }
        }
    }
    if err := models.DeleteUser(id); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusNoContent)
}

func FindUserSubscription(c *fiber.Ctx) error {
    id := c.Params("id")
    isBase64 := false
    if c.Query("base64") == "true" {
        isBase64 = true
    }
    user, err := models.FindUser(id)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.ErrNotFound
        }
        return err
    }
    var subscription strings.Builder
    for _, profile := range user.Profiles {
        subscription.WriteString(profile + "\n")
    }
    if isBase64 {
        return c.SendString(base64.StdEncoding.EncodeToString([]byte(subscription.String())))
    }
    return c.SendString(subscription.String())
}
