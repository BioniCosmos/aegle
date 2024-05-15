package handlers

import (
	"github.com/bionicosmos/aegle/handlers/transfer"
	"github.com/bionicosmos/aegle/models"
	"github.com/bionicosmos/aegle/services"
	"github.com/gofiber/fiber/v2"
)

func FindProfiles(c *fiber.Ctx) error {
	profileNames, err := models.FindProfiles()
	if err != nil {
		return err
	}
	return c.JSON(profileNames)
}

func InsertProfile(c *fiber.Ctx) error {
	body := transfer.InsertProfileBody{}
	if err := c.BodyParser(&body); err != nil {
		return &ParseError{err}
	}
	if err := services.InsertProfile(body.ToProfile()); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func DeleteProfile(c *fiber.Ctx) error {
	if err := services.DeleteProfile(c.Params("name")); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}
