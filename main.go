package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bionicosmos/aegle/config"
	"github.com/bionicosmos/aegle/cron"
	"github.com/bionicosmos/aegle/edge"
	"github.com/bionicosmos/aegle/handlers"
	"github.com/bionicosmos/aegle/models"
	"github.com/bionicosmos/aegle/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`Commands:
  center
  edge`)
		return
	}
	if os.Args[1] == "edge" {
		edge.Start()
	}
	config.Init()
	models.Init()
	cron.Init()

	app := fiber.New()
	app.Use(logger.New())
	routes.Init(app)
	handlers.SessionInit()
	log.Fatal(app.Listen(config.Conf.Listen))
}
