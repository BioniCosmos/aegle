package handler

import (
	"path"

	"github.com/bionicosmos/aegle/config"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	app.Static("/", config.C.Home)
	app.Static("/dashboard/", config.C.Dashboard)
	app.Use("/dashboard/", func(c *fiber.Ctx) error {
		return c.SendFile(path.Join(config.C.Dashboard, "index.html"))
	})

	app.Use("/api", AuthorizeAccount)

	app.Get("/api/node/:id", FindNode)
	app.Get("/api/nodes", FindNodes)
	app.Post("/api/node", InsertNode)
	app.Put("/api/node", UpdateNode)
	app.Delete("/api/node/:id", DeleteNode)

	app.Get("/api/profiles", FindProfiles)
	app.Post("/api/profile", InsertProfile)
	app.Delete("/api/profile/:name", DeleteProfile)

	app.Get("/api/user/:id", FindUser)
	app.Get("/api/users", FindUsers)
	app.Get("/api/user/:id/sub", FindUserProfiles)
	app.Post("/api/user", InsertUser)
	app.Put("/api/user", UpdateUserDate)
	app.Patch("/api/user", UpdateUserProfiles)
	app.Delete("/api/user/:id", DeleteUser)

	app.Post("/api/account/sign-in", SignInAccount)
	app.Post("/api/account/sign-up", SignUpAccount)
	app.Post("/api/account/sign-out", SignOutAccount)
	app.Patch("/api/account/password", ChangeAccountPassword)
	app.Delete("/api/account", DeleteAccount)
}
