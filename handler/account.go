package handler

import (
	"slices"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

func SignUp(c *fiber.Ctx) error {
	body := transfer.SignUpBody{}
	if err := c.BodyParser(&body); err != nil {
		return &ParseError{err}
	}
	if err := service.SignUp(&body); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	account, err := service.Verify(id)
	if err != nil {
		return err
	}
	email := account.Email
	if email == "" {
		return fiber.ErrBadRequest
	}
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	session.Set("email", email)
	if err := session.Save(); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
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
			"/api/account/verification",
			"/api/account/sign-in",
			"/api/user/profiles",
		},
		c.Path(),
	) {
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
