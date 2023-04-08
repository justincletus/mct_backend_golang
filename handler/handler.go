package handler

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/map-backend/database"
	"github.com/justincletus/map-backend/models"
)

//var SECRET = config.SECRET

func GetLocations(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("jwt")
	if cookie != "" {
		var location []models.Location
		_ = database.DB.Select("id", "latitude", "longitude", "city", "state", "country").Find(&location)

		return ctx.Status(200).JSON(fiber.Map{
			"data": location,
		})

	}

	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "user is not authorized",
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
	cookie := ctx.Cookies("jwt")
	if cookie != "" {

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
			address := tempNew["address"].(string)
			if address != "" {
				addressArr := strings.Split(address, ",")
				if addressArr[len(addressArr)-1] != "" {
					location.Country = addressArr[len(addressArr)-1]
				}
				if addressArr[len(addressArr)-2] != "" {
					stateArr := strings.Split(addressArr[len(addressArr)-2], " ")
					location.State = stateArr[1]
					location.PinCode = stateArr[2]
				}

				location.City = addressArr[len(addressArr)-4]
			}

		} else {
			//fmt.Println("no body field")
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

	return ctx.Status(400).JSON(fiber.Map{
		"message": "user is not authorized to perform this task",
	})

}

func GetUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if cookie != "" {
		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET), nil
		})
		if err != nil {
			c.Status(401).JSON(fiber.Map{
				"message": "you are not authorized to perform this task",
			})
		}

		claims := token.Claims.(*jwt.StandardClaims)
		var user models.User
		database.DB.Where("id=?", claims.Issuer).First(&user)
		return c.Status(200).JSON(fiber.Map{
			"id":       user.ID,
			"username": user.Fullname,
		})
	}

	return c.Status(404).JSON(fiber.Map{
		"message": "user is not authorized",
	})
}
