package graphql

import (
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

type Input struct {
	Query         string         `query:"query"`
	OperationName string         `query:"operationName"`
	Variables     map[string]any `query:"variables"`
}

func Handler(c *fiber.Ctx) error {
	var input Input
	if c.Method() == "GET" {
		if err := c.QueryParser(&input); err != nil {
			return c.
				Status(fiber.StatusBadRequest).
				SendString("Cannot parse query parameters: " + err.Error())
		}
	} else if c.Method() == "POST" {
		if err := c.BodyParser(&input); err != nil {
			return c.
				Status(fiber.StatusBadRequest).
				SendString("Cannot parse body: " + err.Error())
		}
	} else {
		return fiber.ErrMethodNotAllowed
	}

	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  input.Query,
		OperationName:  input.OperationName,
		VariableValues: input.Variables,
	})
	c.Set("Content-Type", "application/graphql-response+json")
	return c.JSON(result)
}
