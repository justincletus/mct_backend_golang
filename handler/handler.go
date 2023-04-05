package handler

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	var location models.Location
	var body string

	ctx.BodyParser(&data)

	body = data["body"]
	if body != "" {
		tempNew := make(map[string]interface{})
		json.Unmarshal([]byte(body), &tempNew)
		location.Latitude = strconv.FormatFloat(tempNew["latitude"].(float64), 'f', 2, 64)
		location.Longitude = strconv.FormatFloat(tempNew["longitude"].(float64), 'f', 2, 64)

	} else {
		fmt.Println("no body field")
		location.Latitude = data["latitude"]
		location.Longitude = data["longitude"]

	}

	if location.Latitude != "" && location.Longitude != "" {
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

	return ctx.Status(400).JSON(fiber.Map{
		"message": "latitude or longitude details should be exmpty",
		"data":    data,
	})

}
