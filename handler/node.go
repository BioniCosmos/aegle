package handler

import (
	"context"
	"errors"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindNode(c *fiber.Ctx) error {
	params := transfer.ParamsId{}
	if err := c.ParamsParser(&params); err != nil {
		return &ParseError{err}
	}
	node, err := model.FindNode(&params.Id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.ErrNotFound
		}
		return err
	}
	return c.JSON(node)
}

func FindNodes(c *fiber.Ctx) error {
	pagination, err := service.FindNodes(c.QueryInt("page", 1))
	if err != nil {
		return err
	}
	return c.JSON(pagination)
}

func InsertNode(c *fiber.Ctx) error {
	node := model.Node{}
	if err := c.BodyParser(&node); err != nil {
		return &ParseError{err}
	}
	if err := model.InsertNode(context.Background(), &node); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func UpdateNode(c *fiber.Ctx) error {
	node := model.Node{}
	if err := c.BodyParser(&node); err != nil {
		return &ParseError{err}
	}
	if _, err := model.UpdateNode(
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
	if err := service.DeleteNode(&params.Id); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}
