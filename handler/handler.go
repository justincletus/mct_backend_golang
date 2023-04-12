package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/map-backend/database"
	"github.com/justincletus/map-backend/models"
)

//var SECRET = config.SECRET

func GetLocations(ctx *fiber.Ctx) error {

	var output []models.Location
	_ = database.DB.Preload("Location").Order("created_at desc").Find(&output).Error
	// fmt.Println(temp)
	// fmt.Println(output)

	// cookie := ctx.Cookies("jwt")
	// if cookie != "" {
	// 	var location []models.Location
	// 	_ = database.DB.Select("id", "latitude", "longitude", "city", "state", "country").Find(&location)

	// 	return ctx.Status(200).JSON(fiber.Map{
	// 		"data": location,
	// 	})

	// }

	return ctx.Status(200).JSON(fiber.Map{
		"data": output,
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

func GetLocationByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	var user models.User
	database.DB.Where("username", username).First(&user)
	var locations []models.Location

	database.DB.Where("uid=?", user.Id).Find(&locations)

	return c.Status(200).JSON(fiber.Map{
		"data": locations,
	})
}

func SaveLocation(ctx *fiber.Ctx) error {
	// cookie := ctx.Cookies("jwt")
	// if cookie != "" {

	var data map[string]string
	var location models.Location
	var user models.User
	ctx.BodyParser(&data)

	database.DB.Where("username", data["username"]).Find(&user)
	if user.Id != 0 {
		location.Latitude = data["latitude"]
		location.Longitude = data["longitude"]

		address := data["address"]
		if address != "" {
			addressArr := strings.Split(data["address"], ",")
			if addressArr[len(addressArr)-1] != "" {
				location.Country = addressArr[len(addressArr)-1]
			}
			if addressArr[len(addressArr)-2] != "" {
				stateArr := strings.Split(addressArr[len(addressArr)-2], " ")
				location.State = stateArr[1]
				location.PinCode = stateArr[2]
			}

			location.City = addressArr[len(addressArr)-4]
		} else {
			location.City = data["city"]
			location.State = data["state"]
			location.Country = data["country"]
		}

		location.Uid = user.ID
		location.UserFullName = user.Fullname
		location.User = user

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
	}

	return ctx.Status(400).JSON(fiber.Map{
		"message": "user is not found",
		"data":    data,
	})

	//}

	// return ctx.Status(400).JSON(fiber.Map{
	// 	"message": "user is not authorized to perform this task",
	// })

}
