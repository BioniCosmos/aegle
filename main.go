package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/bionicosmos/aegle/config"
	"github.com/bionicosmos/aegle/edge"
	"github.com/bionicosmos/aegle/handlers"
	"github.com/bionicosmos/aegle/models"
	"github.com/bionicosmos/aegle/routes"
	"github.com/bionicosmos/aegle/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if len(os.Args) < 2 || !slices.Contains([]string{"center", "edge"}, os.Args[1]) {
		fmt.Println(`Commands:
    center
    edge`)
		os.Exit(2)
	}
	params := flag.NewFlagSet("Aegle", flag.ExitOnError)
	configPath := params.String("config", "./config.json", "path of the configuration file")
	params.Parse(os.Args[2:])
	config.Init(*configPath)
	if os.Args[1] == "edge" {
		edge.Start()
	}
	services.Init(models.Init())
	if err := services.CheckUser(); err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			if err := services.CheckUser(); err != nil {
				log.Print(err)
			}
		}
	}()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			var parseError *handlers.ParseError
			if errors.As(err, &parseError) {
				code = fiber.StatusBadRequest
			}
			return c.
				Status(code).
				JSON(map[string]string{"message": err.Error()})
		},
	})
	app.Use(logger.New())
	routes.Init(app)
	handlers.SessionInit()
	log.Fatal(app.Listen(config.C.Listen))
}
