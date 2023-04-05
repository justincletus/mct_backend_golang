package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/map-backend/handler"
)

func Setup(app *fiber.App) {
	app.Get("/", handler.GetLocations)
	app.Get("/:id", handler.GetLocationById)
	app.Post("/", handler.SaveLocation)
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)
	app.Post("/logout", handler.Logout)
}
