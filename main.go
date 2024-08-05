package main

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"time"

	"github.com/bionicosmos/aegle/edge"
	"github.com/bionicosmos/aegle/handler"
	"github.com/bionicosmos/aegle/handler/transfer"
	"github.com/bionicosmos/aegle/model"
	"github.com/bionicosmos/aegle/service"
	"github.com/bionicosmos/aegle/setting"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	mode := loadMode()
	loadEnv(mode)
	if mode == "edge" {
		edge.Start()
		return
	}

	client := model.Init()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()
	service.Init(client)
	setting.Init()

	if err := service.CheckUser(); err != nil {
		panic(err)
	}
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			if err := service.CheckUser(); err != nil {
				log.Print(err)
			}
		}
	}()

	gob.Register(transfer.AccountStatus(""))

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			var parseError *handler.ParseError
			if errors.As(err, &parseError) {
				code = fiber.StatusBadRequest
			}
			return c.
				Status(code).
				JSON(map[string]string{"message": err.Error()})
		},
	})
	app.Use(logger.New())
	handler.Init(app)
	go func() {
		if err := app.Listen(os.Getenv("LISTEN")); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	if err := app.Shutdown(); err != nil {
		panic(err)
	}
}

func loadMode() string {
	if len(os.Args) < 2 ||
		!slices.Contains([]string{"center", "edge"}, os.Args[1]) {
		fmt.Println(`Commands:
    center
    edge`)
		os.Exit(2)
	}
	return os.Args[1]
}

func loadEnv(mode string) {
	if mode == "center" {
		if os.Getenv("DB_URL") == "" {
			log.Fatal("empty DB_URL")
		}
		if os.Getenv("DB_NAME") == "" {
			os.Setenv("DB_NAME", "aegle")
		}
	} else {
		if os.Getenv("XRAY_CONFIG_DIR") == "" {
			os.Setenv("XRAY_CONFIG_DIR", "/usr/local/etc/xray/")
		}
	}
}
