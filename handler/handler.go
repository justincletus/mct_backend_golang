package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/map-backend/database"
	"github.com/justincletus/map-backend/models"
)

func GetUsers(ctx *fiber.Ctx) error {

	return ctx.JSON(fiber.Map{
		"message": "hello world!",
		"status":  fiber.StatusOK,
	})

}

func SaveLocation(ctx *fiber.Ctx) error {
	var data map[string]string

	ctx.BodyParser(&data)

	var location models.Location
	location.Latitude = data["latitude"]
	location.Longtude = data["longtude"]
	err := database.DB.Create(&location).Error
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "location information not saved",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "user location saved",
		"status":  fiber.StatusCreated,
	})

}
