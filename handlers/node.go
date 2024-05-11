package handlers

import (
	"errors"

	"github.com/bionicosmos/aegle/api"
	"github.com/bionicosmos/aegle/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindNode(c *fiber.Ctx) error {
	node, err := models.FindNode(c.Params("id"))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.ErrNotFound
		}
		return err
	}
	return c.JSON(node)
}

func FindNodes(c *fiber.Ctx) error {
	var query struct {
		Skip  int64 `query:"skip"`
		Limit int64 `query:"limit"`
	}
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := models.NodeProfilesAggregateQuery(query.Skip, query.Limit)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func InsertNode(c *fiber.Ctx) error {
	var node models.Node
	if err := c.BodyParser(&node); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := node.Insert(); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusCreated))
}

func UpdateNode(c *fiber.Ctx) error {
	var node models.Node
	if err := c.BodyParser(&node); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := node.Update(); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}

func DeleteNode(c *fiber.Ctx) error {
	id := c.Params("id")
	nodeId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	// Find profiles whose nodeId equals to the id from parameters.
	// If there are errors, return them.
	// If the profile is not empty, also return an error.
	profiles, err := models.FindProfiles(bson.D{
		{Key: "nodeId", Value: nodeId},
	}, bson.D{}, 0, 0)
	if err != nil {
		return err
	}
	if profiles != nil {
		return fiber.NewError(fiber.StatusConflict, "Profiles binding to the node are not empty.")
	}
	if err := models.DeleteNode(c.Params("id")); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}

func ResetNode(c *fiber.Ctx) error {
	id := c.Params("id")
	node, err := models.FindNode(id)
	if err != nil {
		return err
	}
	if err := api.ResetNode(node); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}
