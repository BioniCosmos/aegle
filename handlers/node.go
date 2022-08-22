package handlers

import (
    "github.com/bionicosmos/submgr/models"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

func FindNode(c *fiber.Ctx) error {
    node, err := models.FindNode(c.Params("id"))
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return fiber.ErrNotFound
        }
        return err
    }
    return c.JSON(node)
}

func FindNodes(c *fiber.Ctx) error {
    var query struct {
        skip  int64
        limit int64
    }
    if err := c.QueryParser(&query); err != nil {
        return err
    }

    nodes, err := models.FindNodes(query.skip, query.limit)
    if err != nil {
        return err
    }
    if nodes == nil {
        return fiber.ErrNotFound
    }
    return c.JSON(nodes)
}

func InsertNode(c *fiber.Ctx) error {
    var node models.Node
    if err := c.BodyParser(&node); err != nil {
        return err
    }
    if err := node.Insert(); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusCreated)
}

func UpdateNode(c *fiber.Ctx) error {
    var node models.Node
    if err := c.BodyParser(&node); err != nil {
        return err
    }
    if err := node.Update(c.Params("id")); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusNoContent)
}

func DeleteNode(c *fiber.Ctx) error {
    if err := models.DeleteNode(c.Params("id")); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusNoContent)
}
