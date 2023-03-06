package handlers

import (
	"errors"
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
		if errors.Is(err, mongo.ErrNoDocuments) {
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
		UserId primitive.ObjectID `query:"userId"`
	}
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var nodeIdFilter bson.E
	if !query.NodeId.IsZero() {
		nodeIdFilter = bson.E{Key: "nodeId", Value: query.NodeId}
	}
	var userIdFilter bson.E
	if !query.UserId.IsZero() {
		user, err := models.FindUser(query.UserId.Hex())
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return fiber.NewError(fiber.StatusNotFound, "user not found")
			}
			return err
		}
		profileIds := make(bson.A, 0, len(user.Profiles))
		for profileId := range user.Profiles {
			profileIds = append(profileIds, profileId)
		}
		userIdFilter = bson.E{Key: "_id", Value: bson.D{
			{Key: "$in", Value: profileIds},
		}}
	}
	profiles, err := models.FindProfiles(bson.D{
		nodeIdFilter,
		userIdFilter,
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "node not found")
		}
		return err
	}
	profile.Inbound.Tag = profile.Id.Hex()
	if err := api.AddInbound(profile.Inbound.ToConf(), node.APIAddress); err != nil {
		return err
	}
	if err := profile.Insert(); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusCreated))
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
		return fiber.NewError(fiber.StatusConflict, "Users binding to the profile are not empty.")
	}

	profile, err := models.FindProfile(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(fiber.StatusNotFound, "profile not found")
		}
		return err
	}

	node, err := models.FindNode(profile.NodeId)
	if err != nil {
		return err
	}
	if err := api.RemoveInbound(profile.Inbound.Tag, node.APIAddress); err != nil {
		return err
	}
	if err := models.DeleteProfile(id); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}
