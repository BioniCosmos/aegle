package main

import (
	"log"

	"github.com/bionicosmos/aegle/config"
	"github.com/bionicosmos/aegle/cron"
	"github.com/bionicosmos/aegle/handlers"
	"github.com/bionicosmos/aegle/models"
	"github.com/bionicosmos/aegle/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.Init()
	models.Init()
	cron.Init()

	app := fiber.New()
	app.Use(logger.New())
	routes.Init(app)
	handlers.SessionInit()
	log.Fatal(app.Listen(config.Conf.Listen))
}
