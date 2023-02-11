package routes

import (
    "path"

    "github.com/bionicosmos/submgr/config"
    "github.com/bionicosmos/submgr/handlers"
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

    app.Get("/api/profile/:id", handlers.FindProfile)
    app.Get("/api/profiles", handlers.FindProfiles)
    app.Post("/api/profile", handlers.InsertProfile)
    app.Delete("/api/profile/:id", handlers.DeleteProfile)

    app.Get("/api/user/:id", handlers.FindUser)
    app.Get("/api/users", handlers.FindUsers)
    app.Post("/api/user", handlers.InsertUser)
    app.Patch("/api/user/:id", handlers.UpdateUser)
    app.Delete("/api/user/:id", handlers.DeleteUser)
    app.Get("/api/user/:id/sub", handlers.FindUserSubscription)

    app.Get("/api/invoice/:id", handlers.FindUserInvoice)
    app.Get("/api/invoices", handlers.FindUserInvoices)
    app.Patch("/api/invoice/:id", handlers.ExtendBillingDate)

    app.Post("/api/account/sign-in", handlers.SignInAccount)
    app.Post("/api/account/sign-up", handlers.SignUpAccount)
    app.Patch("/api/account/password", handlers.ChangeAccountPassword)
    app.Delete("/api/account", handlers.DeleteAccount)
}
