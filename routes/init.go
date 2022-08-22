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

    app.Get("/api/inbound/:id", handlers.FindInbound)
    app.Get("/api/inbounds", handlers.FindInbounds)
    app.Post("/api/inbound", handlers.InsertInbound)
    app.Delete("/api/inbound/:id", handlers.DeleteInbound)
}
