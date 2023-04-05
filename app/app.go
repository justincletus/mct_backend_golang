package app

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/justincletus/map-backend/database"
	"github.com/justincletus/map-backend/router"
)

func Start() {
	project := fiber.New()
	if err := database.Connetion(); err != nil {
		log.Fatal(err)
	}

	project.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	router.Setup(project)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(project.Listen(":" + port))

}
