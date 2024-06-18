package handler

import (
	"path"
	"time"

	"github.com/bionicosmos/aegle/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mongodb"
)

func Init(app *fiber.App) {
	app.Static("/", config.C.Home)
	app.Static("/dashboard/", config.C.Dashboard)
	app.Use("/dashboard/", func(c *fiber.Ctx) error {
		return c.SendFile(path.Join(config.C.Dashboard, "index.html"))
	})

	app.Use("/api", Auth)

	app.Get("/api/node/:id", FindNode)
	app.Get("/api/nodes", FindNodes)
	app.Post("/api/node", InsertNode)
	app.Put("/api/node", UpdateNode)
	app.Delete("/api/node/:id", DeleteNode)

	app.Get("/api/profiles", FindProfiles)
	app.Post("/api/profile", InsertProfile)
	app.Delete("/api/profile/:name", DeleteProfile)

	app.Get("/api/user/profiles", FindUserProfiles)
	app.Get("/api/user/:id", FindUser)
	app.Get("/api/users", FindUsers)
	app.Post("/api/user", InsertUser)
	app.Put("/api/user", UpdateUserDate)
	app.Patch("/api/user", UpdateUserProfiles)
	app.Delete("/api/user/:id", DeleteUser)

	app.Post("/api/account/sign-up", SignUp)
	app.Post("/api/account/verification", Verify)
	app.Post("/api/account/sign-in", SignIn)

	store = session.New(session.Config{
		Expiration: time.Hour * 24 * 365,
		Storage: mongodb.New(mongodb.Config{
			ConnectionURI: config.C.DatabaseURL,
			Database:      config.C.DatabaseName,
		}),
	})
}
