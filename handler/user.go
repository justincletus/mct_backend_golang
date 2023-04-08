package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/map-backend/config"
	"github.com/justincletus/map-backend/database"
	"github.com/justincletus/map-backend/models"
	"golang.org/x/crypto/bcrypt"
)

var SECRET = config.SECRET

func Register(c *fiber.Ctx) error {
	var data map[string]string

	var user models.User

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "user registration failed",
		})
	}

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error in creating hashed password",
		})
	}

	user.Email = data["email"]
	user.Fullname = data["fullname"]
	user.Mobile = data["mobile"]
	user.Password = passwordBytes
	database.DB.Create(&user)

	return c.Status(201).JSON(fiber.Map{
		"response": user,
	})

}

func Login(c *fiber.Ctx) error {

	var data map[string]string

	var user models.User

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "something went wrong please check the username/email and password",
		})
	}
	database.DB.Where("email", data["email"]).First(&user)

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "user authentication failed",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SECRET))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "unable decrypt the password, please check the password again",
		})

	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.Status(200).JSON(fiber.Map{
		"message": "login successful!",
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	fmt.Println(cookie)

	return c.Status(200).JSON(fiber.Map{
		"message": "logout suuccess!",
	})

}
