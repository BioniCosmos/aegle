package routes

import (
    "github.com/bionicosmos/submgr/models"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo"
)

func Init(app *fiber.App) {
    app.Get("/api/node/:id", func(c *fiber.Ctx) error {
        node, err := models.FindNode(c.Params("id"))
        if err != nil {
            if err == mongo.ErrNoDocuments {
                return fiber.ErrNotFound
            }
            return err
        }
        return c.JSON(node)
    })
    app.Get("/api/nodes", func(c *fiber.Ctx) error {
        var query struct {
            skip  int64
            limit int64
        }
        err := c.QueryParser(&query)
        if err != nil {
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
    })
    app.Post("/api/node", func(c *fiber.Ctx) error {
        var node models.Node
        err := c.BodyParser(&node)
        if err != nil {
            return err
        }
        err = node.Insert()
        if err != nil {
            return err
        }
        return c.SendStatus(fiber.StatusCreated)
    })
    app.Patch("/api/node/:id", func(c *fiber.Ctx) error {
        var node models.Node
        err := c.BodyParser(&node)
        if err != nil {
            return err
        }
        err = node.Update(c.Params("id"))
        if err != nil {
            return err
        }
        return c.SendStatus(fiber.StatusNoContent)
    })
    app.Delete("/api/node/:id", func(c *fiber.Ctx) error {
        err := models.DeleteNode(c.Params("id"))
        if err != nil {
            return err
        }
        return c.SendStatus(fiber.StatusNoContent)
    })
}
