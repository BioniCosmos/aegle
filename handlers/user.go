package handlers

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/bionicosmos/submgr/api"
	"github.com/bionicosmos/submgr/handlers/transfer"
	"github.com/bionicosmos/submgr/models"
	"github.com/bionicosmos/submgr/services"
	"github.com/bionicosmos/submgr/services/subscription"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUser(c *fiber.Ctx) error {
	user, err := models.FindUser(c.Params("id"))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.ErrNotFound
		}
		return err
	}
	return c.JSON(transfer.UserGetFromUser(user))
}

func FindUsers(c *fiber.Ctx) error {
	var query struct {
		Skip  int64 `query:"skip"`
		Limit int64 `query:"limit"`
	}
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	res, err := models.UserProfilesAggregateQuery(query.Skip, query.Limit)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func InsertUser(c *fiber.Ctx) error {
	userPost := transfer.UserPost{}
	if err := c.BodyParser(&userPost); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	user := userPost.ToUser()
	for _, profileId := range userPost.ProfileIds {
		profile, err := models.FindProfile(profileId.Hex())
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return fiber.NewError(fiber.StatusUnprocessableEntity, "profile not found")
			}
			return err
		}

		node, err := models.FindNode(profile.NodeId)
		if err != nil {
			return err
		}
		if err := api.AddUser(
			profile.Inbound.ToConf(),
			&user,
			node.APIAddress,
		); err != nil {
			return err
		}
		user.Profiles[profileId], err = subscription.Generate(profile, user.Account)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}
	if err := user.Insert(); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusCreated))
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	type operation uint
	const (
		Add operation = iota
		Remove
	)
	var body struct {
		Id        primitive.ObjectID `json:"id"`
		Operation operation          `json:"operation"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	user, err := models.FindUser(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return err
	}
	profile, err := models.FindProfile(body.Id.Hex())
	if err != nil {
		return err
	}
	if body.Operation == Add {
		if err := services.UserAddProfile(user, profile); err != nil {
			subscriptionError := new(subscription.Error)
			if errors.As(err, &subscriptionError) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return err
		}
	} else {
		if err := services.UserRemoveProfile(user, profile); err != nil {
			return err
		}
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := models.FindUser(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return err
	}
	for profileId := range user.Profiles {
		profile, err := models.FindProfile(profileId.Hex())
		if err != nil {
			return err
		}

		node, err := models.FindNode(profile.NodeId)
		if err != nil {
			return err
		}
		if err := api.RemoveUser(profile.Inbound.ToConf(), user, node.APIAddress); err != nil {
			return err
		}
	}
	if err := models.DeleteUser(id); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}

func FindUserSubscription(c *fiber.Ctx) error {
	id := c.Params("id")
	isBase64 := false
	if c.Query("base64") == "true" {
		isBase64 = true
	}
	user, err := models.FindUser(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.ErrNotFound
		}
		return err
	}
	var subscription strings.Builder
	for _, profile := range user.Profiles {
		subscription.WriteString(profile + "\n")
	}
	if isBase64 {
		return c.SendString(base64.StdEncoding.EncodeToString([]byte(subscription.String())))
	}
	return c.SendString(subscription.String())
}

func Extend(c *fiber.Ctx) error {
	userPut := transfer.UserPut{}
	if err := c.BodyParser(&userPut); err != nil {
		return fiber.ErrBadRequest
	}
	user, err := models.FindUser(userPut.Id)
	if err != nil {
		return err
	}
	user.Cycles = userPut.Cycles
	user.NextDate = userPut.NextDate
	if err := user.Update(); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
