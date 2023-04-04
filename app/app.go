package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/gpstrack/router"
)

func Start() {
	project := fiber.New()

	router.Setup(project)

	log.Fatal(project.Listen(":8000"))

}
