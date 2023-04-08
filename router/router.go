package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/justincletus/map-backend/handler"
)

func Setup(app *fiber.App) {
	api := app.Group("/api", logger.New())

	api.Get("/", handler.GetLocations)
	api.Post("/", handler.SaveLocation)
	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)
	api.Post("/logout", handler.Logout)
	api.Get("/user", handler.GetUser)
	api.Get("/location/:id", handler.GetLocationById)
}
