package handler

import (
	"errors"
	"slices"
	"strings"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/model/account"
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
	account, err := service.SignUp(&body)
	if err != nil {
		if errors.Is(err, service.ErrAccountExists) {
			return fiber.NewError(fiber.StatusConflict, "The email exists.")
		}
		return err
	}
	if err := setAccount(c, &account); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(account)
}

func Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	account := getAccount(c)
	account, err := service.Verify(id, &account)
	if err != nil {
		if errors.Is(err, service.ErrVerified) {
			return fiber.ErrConflict
		}
		if errors.Is(err, service.ErrLinkExpired) {
			return fiber.ErrNotFound
		}
		return err
	}
	if err := setAccount(c, &account); err != nil {
		return err
	}
	return c.JSON(transfer.AccountOmitPassword(&account))
}

func SendVerificationLink(c *fiber.Ctx) error {
	account := getAccount(c)
	if err := service.SendVerificationLink(&account); err != nil {
		if errors.Is(err, service.ErrVerified) {
			return fiber.ErrConflict
		}
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func SignIn(c *fiber.Ctx) error {
	body := transfer.SignInBody{}
	if err := c.BodyParser(&body); err != nil {
		return &ParseError{err}
	}
	account, err := service.SignIn(&body)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(fiber.StatusNotFound, "user does not exist")
		}
		if errors.Is(err, service.ErrPassword) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return err
	}
	if err := setAccount(c, &account); err != nil {
		return err
	}
	return c.JSON(transfer.AccountOmitPassword(&account))
}

func Auth(c *fiber.Ctx) error {
	if slices.Contains(
		[]string{
			"/api/account/sign-up",
			"/api/account/sign-in",
			"/api/user/profiles",
		},
		c.Path(),
	) {
		return c.Next()
	}
	account, err := useAccount(c)
	if err != nil {
		return err
	}
	if !access(&account, c.Path()) {
		return fiber.ErrForbidden
	}
	c.Locals("account", account)
	return c.Next()
}

func useAccount(c *fiber.Ctx) (model.Account, error) {
	session, err := store.Get(c)
	if err != nil {
		return model.Account{}, err
	}
	account, ok := session.Get("account").(model.Account)
	if !ok {
		return model.Account{}, fiber.ErrUnauthorized
	}
	return account, nil
}

func getAccount(c *fiber.Ctx) model.Account {
	return c.Locals("account").(model.Account)
}

func setAccount(c *fiber.Ctx, account *model.Account) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	session.Set("account", *account)
	return session.Save()
}

func access(a *model.Account, path string) bool {
	return strings.HasPrefix(path, "/api/account/verification") ||
		(a.Status == account.Normal && a.Role == account.Admin)
}
