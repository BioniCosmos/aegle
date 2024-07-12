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

func GetAccount(c *fiber.Ctx) error {
	return c.JSON(transfer.ToAccount(getSession(c)))
}

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
	if err := setSession(c, &account, transfer.AccountUnverified); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	account, _ := getSession(c)
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
	if err := setSession(c, &account, transfer.AccountSignedIn); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func SendVerificationLink(c *fiber.Ctx) error {
	account, _ := getSession(c)
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
	status := transfer.AccountSignedIn
	if account.TOTP != "" {
		status = transfer.AccountNeedMFA
	}
	if err := setSession(c, &account, status); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func CreateTOTP(c *fiber.Ctx) error {
	account, _ := getSession(c)
	body, err := service.CreateTOTP(account.Email)
	if err != nil {
		return err
	}
	return c.JSON(body)
}

func ConfirmTOTP(c *fiber.Ctx) error {
	body := transfer.ConfirmTOTPBody{}
	if err := c.BodyParser(&body); err != nil {
		return &ParseError{err}
	}
	account, _ := getSession(c)
	if err := service.ConfirmTOTP(account.Email, body.Code); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.NewError(
				fiber.StatusBadRequest,
				"Uninitialized or expired TOTP",
			)
		}
		if errors.Is(err, service.ErrInvalidTOTP) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid TOTP")
		}
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func DeleteTOTP(c *fiber.Ctx) error {
	account, _ := getSession(c)
	if err := service.DeleteTOTP(account.Email); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fiber.ErrNotFound
		}
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
	) {
		return c.Next()
	}
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	if session.Fresh() {
		return fiber.ErrUnauthorized

	}
	account := session.Get("account").(model.Account)
	status := session.Get("status").(transfer.AccountStatus)
	if !access(&account, status, c.Path()) {
		return fiber.ErrForbidden
	}
	c.Locals("account", account)
	c.Locals("status", status)
	return c.Next()
}

func getSession(c *fiber.Ctx) (model.Account, transfer.AccountStatus) {
	return c.Locals("account").(model.Account),
		c.Locals("status").(transfer.AccountStatus)
}

func setSession(
	c *fiber.Ctx,
	account *model.Account,
	status transfer.AccountStatus,
) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	session.Set("account", *account)
	session.Set("status", status)
	return session.Save()
}

func access(a *model.Account, status transfer.AccountStatus, path string) bool {
	return strings.HasPrefix(path, "/api/account/verification") ||
		(a.Status == account.Normal &&
			a.Role == account.Admin &&
			status == transfer.AccountSignedIn)
}
