package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/map-backend/handler"
)

func Setup(app *fiber.App) {
	app.Get("/", handler.GetUsers)
}
