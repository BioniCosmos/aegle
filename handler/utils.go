package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func toJSON(c *fiber.Ctx, status int) error {
	return c.
		Status(status).
		JSON(map[string]string{"message": utils.StatusMessage(status)})
}
