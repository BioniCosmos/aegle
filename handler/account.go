package handler

import (
	"slices"
	"strings"

	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const (
	sessionEmail  = "email"
	sessionStatus = "status"
)

var store *session.Store

func GetAccount(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	if session.Fresh() {
		return c.JSON(nil)
	}
	email := session.Get(sessionEmail).(string)
	status := session.Get(sessionStatus).(transfer.AccountStatus)
	account, err := model.FindAccount(email)
	if err != nil {
		return err
	}
	return c.JSON(transfer.ToAccount(&account, status))
}

func SignUp(c *fiber.Ctx) error {
	body := transfer.SignUpBody{}
	if err := c.BodyParser(&body); err != nil {
		return &ParseError{err}
	}
	account, err := service.SignUp(&body)
	if err != nil {
		return err
	}
	if err := setSession(
		c,
		account.Email,
		transfer.AccountUnverified,
	); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	email, _ := getSession(c)
	if err := service.Verify(id, email); err != nil {
		return err
	}
	if err := setSession(c, email, transfer.AccountSignedIn); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func SendVerificationLink(c *fiber.Ctx) error {
	email, _ := getSession(c)
	if err := service.SendVerificationLink(email); err != nil {
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
		return err
	}
	status := transfer.AccountSignedIn
	if account.MFA.TOTP != "" {
		status = transfer.AccountNeedMFA
	}
	if err := setSession(c, account.Email, status); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func FindMFA(c *fiber.Ctx) error {
	email, _ := getSession(c)
	account, err := model.FindAccount(email)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"totp": account.MFA.TOTP != ""})
}

func MFA(c *fiber.Ctx) error {
	email, status := getSession(c)
	if status == transfer.AccountSignedIn {
		return fiber.ErrConflict
	}
	body := transfer.MFABody{}
	if err := c.BodyParser(&body); err != nil {
		return &ParseError{err}
	}
	if err := service.MFA(email, &body); err != nil {
		return err
	}
	if err := setSession(c, email, transfer.AccountSignedIn); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func CreateTOTP(c *fiber.Ctx) error {
	email, _ := getSession(c)
	body, err := service.CreateTOTP(email)
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
	email, _ := getSession(c)
	if err := service.ConfirmTOTP(email, body.Code); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusCreated)
}

func DeleteTOTP(c *fiber.Ctx) error {
	email, _ := getSession(c)
	if err := service.DeleteTOTP(email); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func GetRecoveryCodes(c *fiber.Ctx) error {
	email, _ := getSession(c)
	codes, err := service.GetRecoveryCodes(email)
	if err != nil {
		return err
	}
	return c.JSON(codes)
}

func RenewRecoveryCodes(c *fiber.Ctx) error {
	email, _ := getSession(c)
	if _, err := service.RenewRecoveryCodes(email); err != nil {
		return err
	}
	return toJSON(c, fiber.StatusOK)
}

func Auth(c *fiber.Ctx) error {
	if slices.Contains(
		[]string{
			"/api/account",
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
	email := session.Get(sessionEmail).(string)
	status := session.Get(sessionStatus).(transfer.AccountStatus)

	isSignedIn := status == transfer.AccountSignedIn
	toMFA := status == transfer.AccountNeedMFA && c.Path() == "/api/account/mfa"
	toVerify := status == transfer.AccountUnverified &&
		strings.HasPrefix(c.Path(), "/api/account/verification")
	if !isSignedIn && !toMFA && !toVerify {
		return fiber.ErrForbidden
	}

	c.Locals(sessionEmail, email)
	c.Locals(sessionStatus, status)
	return c.Next()
}

func getSession(c *fiber.Ctx) (string, transfer.AccountStatus) {
	return c.Locals(sessionEmail).(string),
		c.Locals(sessionStatus).(transfer.AccountStatus)
}

func setSession(
	c *fiber.Ctx,
	email string,
	status transfer.AccountStatus,
) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}
	session.Set(sessionEmail, email)
	session.Set(sessionStatus, status)
	return session.Save()
}
