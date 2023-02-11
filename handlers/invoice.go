package handlers

import (
	"time"

	"github.com/bionicosmos/submgr/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type invoice struct {
	Id              primitive.ObjectID `json:"id"`
	Name            string             `json:"name"`
	NextBillingDate *time.Time         `json:"nextBillingDate"`
	IsPaid          bool               `json:"isPaid"`
}

func FindUserInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := models.FindUser(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.ErrNotFound
		}
		return err
	}
	return c.JSON(invoiceFromUser(user))
}

func FindUserInvoices(c *fiber.Ctx) error {
	var query struct {
		Skip  int64 `query:"skip"`
		Limit int64 `query:"limit"`
	}
	if err := c.QueryParser(&query); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	users, err := models.FindUsers(bson.D{}, bson.D{
		{Key: "name", Value: 1},
	}, query.Skip, query.Limit)
	if err != nil {
		return err
	}
	if users == nil {
		return fiber.ErrNotFound
	}
	invoices := make([]invoice, 0, len(users))
	for _, user := range users {
		invoices = append(invoices, *invoiceFromUser(&user))
	}
	return c.JSON(invoices)
}

func ExtendBillingDate(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		BillingDate *time.Time `json:"billingDate"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	user, err := models.FindUser(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.ErrNotFound
		}
		return err
	}
	user.BillingDate = body.BillingDate
	if err := user.Update(); err != nil {
		return err
	}
	return c.JSON(fiber.NewError(fiber.StatusOK))
}

func invoiceFromUser(user *models.User) *invoice {
	nextBillingDate := user.BillingDate.AddDate(0, 1, 0)
	isPaid := false
	if nextBillingDate.After(time.Now()) {
		isPaid = true
	}
	return &invoice{
		Id:              user.Id,
		Name:            user.Name,
		NextBillingDate: &nextBillingDate,
		IsPaid:          isPaid,
	}
}
