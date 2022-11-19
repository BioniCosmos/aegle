package routes

import (
    "github.com/bionicosmos/submgr/handlers"
    "github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
    app.Get("/api/node/:id", handlers.FindNode)
    app.Get("/api/nodes", handlers.FindNodes)
    app.Post("/api/node", handlers.InsertNode)
    app.Patch("/api/node/:id", handlers.UpdateNode)
    app.Delete("/api/node/:id", handlers.DeleteNode)

    app.Get("/api/profile/:id", handlers.FindProfile)
    app.Get("/api/profiles", handlers.FindProfiles)
    app.Post("/api/profile", handlers.InsertProfile)
    app.Delete("/api/profile/:id", handlers.DeleteProfile)
}
