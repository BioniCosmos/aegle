package routes

import (
    "github.com/bionicosmos/submgr/handlers"
    "github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
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
}
