package app

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/map-backend/database"
	"github.com/justincletus/map-backend/router"
)

func Start() {
	project := fiber.New()
	if err := database.Connetion(); err != nil {
		log.Fatal(err)
	}

	router.Setup(project)

	log.Fatal(project.Listen(":8000"))

}
