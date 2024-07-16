package handler

import (
	"os"
	"path"
	"time"

	"github.com/bionicosmos/aegle/setting"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/mongodb"
)

func Init(app *fiber.App) {
	app.Static("/", setting.X.Home)
	app.Static("/dashboard/", setting.X.Dashboard)
	app.Use("/dashboard/", func(c *fiber.Ctx) error {
		return c.SendFile(path.Join(setting.X.Dashboard, "index.html"))
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

	app.Get("/api/account", GetAccount)
	app.Post("/api/account/sign-up", SignUp)
	app.Post("/api/account/verification/:id", Verify)
	app.Post("/api/account/verification", SendVerificationLink)
	app.Post("/api/account/sign-in", SignIn)
	app.Get("/api/account/mfa", FindMFA)
	app.Post("/api/account/mfa", MFA)
	app.Get("/api/account/mfa/totp", CreateTOTP)
	app.Post("/api/account/mfa/totp", ConfirmTOTP)
	app.Delete("/api/account/mfa/totp", DeleteTOTP)

	store = session.New(session.Config{
		Expiration: time.Hour * 24 * 365,
		Storage: mongodb.New(mongodb.Config{
			ConnectionURI: os.Getenv("DB_URL"),
			Database:      os.Getenv("DB_NAME"),
		}),
	})
}
