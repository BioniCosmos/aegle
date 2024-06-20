package handler

import (
	"errors"
	"slices"
	"strings"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.mongodb.org/mongo-driver/mongo"
)

var store *session.Store

func SignUp(c *fiber.Ctx) error {
	body := transfer.SignUpBody{}
	if err := c.BodyParser(&body); err != nil {
		return &ParseError{err}
	}
	if err := service.SignUp(&body); err != nil {
		if errors.Is(err, service.ErrAccountExists) {
			return fiber.NewError(fiber.StatusConflict, "The email exists.")
		}
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	account, err := service.Verify(id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.ErrNotFound
		}
		return err
	}
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	session.Set("email", account.Email)
	if err := session.Save(); err != nil {
		return err
	}
	return c.JSON(transfer.AccountOmitPassword(&account))
}

func SignIn(c *fiber.Ctx) error {
	body := transfer.SignInBody{}
	success, err := service.SignIn(&body)
	if err != nil {
		return err
	}
	if !success {
		return fiber.ErrBadRequest
	}
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	session.Set("email", body.Email)
	if err := session.Save(); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func Auth(c *fiber.Ctx) error {
	if slices.Contains(
		[]string{
			"/api/account/sign-up",
			"/api/account/sign-in",
			"/api/user/profiles",
		},
		c.Path(),
	) || strings.HasPrefix(c.Path(), "/api/account/verification") {
		return c.Next()
	}
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	if session.Get("email") == nil {
		return fiber.ErrUnauthorized
	}
	return c.Next()
}
