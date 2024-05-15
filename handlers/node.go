package handlers

import (
	"context"
	"errors"

	"github.com/bionicosmos/aegle/handlers/transfer"
	"github.com/bionicosmos/aegle/models"
	"github.com/bionicosmos/aegle/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindNode(c *fiber.Ctx) error {
	params := transfer.ParamsId{}
	if err := c.ParamsParser(&params); err != nil {
		return &ParseError{err}
	}
	node, err := models.FindNode(&params.Id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.ErrNotFound
		}
		return err
	}
	return c.JSON(node)
}

func FindNodes(c *fiber.Ctx) error {
	query := models.Query{}
	if err := c.QueryParser(&query); err != nil {
		return &ParseError{err}
	}
	nodes, err := models.FindNodes(&query)
	if err != nil {
		return err
	}
	return c.JSON(nodes)
}

func InsertNode(c *fiber.Ctx) error {
	node := models.Node{}
	if err := c.BodyParser(&node); err != nil {
		return &ParseError{err}
	}
	if err := models.InsertNode(context.Background(), &node); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func UpdateNode(c *fiber.Ctx) error {
	node := models.Node{}
	if err := c.BodyParser(&node); err != nil {
		return &ParseError{err}
	}
	if _, err := models.UpdateNode(
		context.Background(),
		bson.M{"_id": node.Id},
		bson.M{"$set": node},
	); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func DeleteNode(c *fiber.Ctx) error {
	params := transfer.ParamsId{}
	if err := c.ParamsParser(&params); err != nil {
		return err
	}
	if err := services.DeleteNode(&params.Id); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}
