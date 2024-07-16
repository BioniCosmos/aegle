package handler

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUser(c *fiber.Ctx) error {
	params := transfer.ParamsId{}
	if err := c.ParamsParser(&params); err != nil {
		return &ParseError{err}
	}
	user, err := model.FindUser(context.Background(), &params.Id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.ErrNotFound
		}
		return err
	}
	return c.JSON(user)
}

func FindUsers(c *fiber.Ctx) error {
	pagination, err := service.FindUsers(c.QueryInt("page", 1))
	if err != nil {
		return err
	}
	return c.JSON(pagination)
}

func FindUserProfiles(c *fiber.Ctx) error {
	query := transfer.FindUserProfilesQuery{}
	if err := c.QueryParser(&query); err != nil {
		return &ParseError{err}
	}
	user, err := model.FindUser(context.Background(), &query.Id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.ErrNotFound
		}
		return err
	}
	links := make([]string, 0)
	for _, profile := range user.Profiles {
		links = append(links, profile.Link)
	}
	s := strings.Join(links, "\n")
	if query.Base64 {
		s = base64.StdEncoding.EncodeToString([]byte(s))
	}
	return c.SendString(s)
}

func InsertUser(c *fiber.Ctx) error {
	user := model.User{}
	if err := c.BodyParser(&user); err != nil {
		return &ParseError{err}
	}
	id, err := model.InsertUser(context.Background(), &user)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func UpdateUserDate(c *fiber.Ctx) error {
	body := transfer.UpdateUserDateBody{}
	if err := c.BodyParser(&body); err != nil {
		return &ParseError{err}
	}
	if err := model.UpdateUser(
		context.Background(),
		bson.M{"_id": body.Id},
		bson.M{
			"$set": bson.M{"cycles": body.Cycles, "nextDate": body.NextDate},
		}); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func UpdateUserProfiles(c *fiber.Ctx) error {
	body := transfer.UpdateUserProfilesBody{}
	if err := c.BodyParser(&body); err != nil {
		return &ParseError{err}
	}
	if err := service.UpdateUserProfiles(&body); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func DeleteUser(c *fiber.Ctx) error {
	params := transfer.ParamsId{}
	if err := c.ParamsParser(&params); err != nil {
		return &ParseError{err}
	}
	if err := service.DeleteUser(&params.Id); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}
