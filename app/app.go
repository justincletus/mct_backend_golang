package app

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/gofiber/fiber/v2/middleware/multipart"
	"github.com/justincletus/mct_backend_golang/database"
	"github.com/justincletus/mct_backend_golang/router"
)

func Start() {
	project := fiber.New()
	if err := database.Connetion(); err != nil {
		log.Fatal(err)
	}

	project.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://127.0.0.1:3000, https://*.megachipstechnology.online, https://megachipstechnology.online",
		AllowHeaders:     "Origin, Content-Type, Accept, role",
		AllowCredentials: true,
	}))
	//project.Use(multipart.New())

	router.Setup(project)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(project.Listen(":" + port))

}
