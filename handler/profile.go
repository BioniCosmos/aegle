package handler

import (
	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/service"
	"github.com/gofiber/fiber/v2"
)

func FindProfiles(c *fiber.Ctx) error {
	profileNames, err := model.FindProfiles()
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
	if err := service.InsertProfile(body.ToProfile()); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func DeleteProfile(c *fiber.Ctx) error {
	if err := service.DeleteProfile(c.Params("name")); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}
