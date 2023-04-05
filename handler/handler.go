package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/map-backend/database"
	"github.com/justincletus/map-backend/models"
)

func GetLocations(ctx *fiber.Ctx) error {

	var location []models.Location
	_ = database.DB.Select("id", "latitude", "longitude").Find(&location)

	return ctx.JSON(fiber.Map{
		"data": location,
	})

}

func GetLocationById(c *fiber.Ctx) error {
	lid := c.Params("id")
	var location models.Location
	database.DB.Where("id", lid).First(&location)
	if location.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "location not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": location,
	})
}

func SaveLocation(ctx *fiber.Ctx) error {
	var data map[string]string

	ctx.BodyParser(&data)

	var location models.Location
	location.Latitude = data["latitude"]
	location.Longitude = data["longtude"]
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
