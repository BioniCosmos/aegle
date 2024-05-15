package routes

import (
	"path"

	"github.com/bionicosmos/aegle/config"
	"github.com/bionicosmos/aegle/handlers"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	app.Static("/", config.Conf.Static.Home)
	app.Static("/dashboard/", config.Conf.Static.Dashboard)
	app.Use("/dashboard/", func(c *fiber.Ctx) error {
		return c.SendFile(path.Join(config.Conf.Static.Dashboard, "index.html"))
	})

	app.Use("/api", handlers.AuthorizeAccount)

	app.Get("/api/node/:id", handlers.FindNode)
	app.Get("/api/nodes", handlers.FindNodes)
	app.Post("/api/node", handlers.InsertNode)
	app.Put("/api/node", handlers.UpdateNode)
	app.Delete("/api/node/:id", handlers.DeleteNode)

	app.Get("/api/profiles", handlers.FindProfiles)
	app.Post("/api/profile", handlers.InsertProfile)
	app.Delete("/api/profile/:name", handlers.DeleteProfile)

	app.Get("/api/user/:id", handlers.FindUser)
	app.Get("/api/users", handlers.FindUsers)
	app.Get("/api/user/:id/sub", handlers.FindUserProfiles)
	app.Post("/api/user", handlers.InsertUser)
	app.Put("/api/user", handlers.UpdateUserDate)
	app.Patch("/api/user", handlers.UpdateUserProfiles)
	app.Delete("/api/user/:id", handlers.DeleteUser)

	app.Post("/api/account/sign-in", handlers.SignInAccount)
	app.Post("/api/account/sign-up", handlers.SignUpAccount)
	app.Post("/api/account/sign-out", handlers.SignOutAccount)
	app.Patch("/api/account/password", handlers.ChangeAccountPassword)
	app.Delete("/api/account", handlers.DeleteAccount)
}
