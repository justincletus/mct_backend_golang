package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/cms/config"
	"github.com/justincletus/cms/database"
	"github.com/justincletus/cms/models"
	"github.com/justincletus/cms/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var SECRET = config.SECRET

func Register(c *fiber.Ctx) error {
	var data map[string]string

	var user models.User
	var manager models.Manager

	err := c.BodyParser(&data)
	if err != nil {
		fmt.Println(err)
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
	user.Username = getUsername(data["fullname"])
	if data["user_role"] == "" {
		data["user_role"] = "user"
	}

	user.Role = data["user_role"]
	user.Code = utils.CreateUuid()

	if data["email_verified"] == "true" {
		user.EmailVerified = true
	} else {
		user.EmailVerified = false
	}

	uId := database.DB.Create(&user)
	rowId := uId.RowsAffected

	if rowId == 0 {
		if !errors.Is(uId.Error, gorm.ErrRecordNotFound) {
			return c.Status(400).JSON(fiber.Map{
				"message": "duplicate email id found, please use new email id",
			})
		} else {
			fmt.Println("something went wrong")
		}
	}

	if data["role"] == "manager" {
		manager.UserId = user.Id
	}

	if manager.UserId != 0 {
		database.DB.Create(&manager)
	}

	if user.EmailVerified {
		emailStruct := utils.EmailBody{
			Email:    user.Email,
			Password: data["password"],
		}
		err = emailStruct.SendEmail(user.Email, "Account Created", "account.html")
		if err != nil {
			fmt.Println(err)
		}

		return c.Status(201).JSON(fiber.Map{
			"data": data,
		})
	} else {

		host := utils.GetRemoteHostAddress(c)

		e_mail_struct := utils.EmailBody{
			Code: user.Code,
			Url:  host,
		}
		err = e_mail_struct.SendEmail(user.Email, "Account Confirmation", "email_verification.html")
		if err != nil {
			fmt.Println(err)
		}

		return c.Status(201).JSON(fiber.Map{
			"data": user,
		})
	}

}

func getUsername(name string) string {
	name = strings.ToLower(name)
	if strings.Contains(name, " ") {
		names := strings.Join(strings.Split(name, " "), "")
		name = names
	}
	if len(name) > 5 {
		name = name[0:4]
	}
	timeValue := strconv.Itoa(int(time.Now().Unix()))
	name += timeValue

	return name
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
	if user.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if !(user.EmailVerified) {
		return c.Status(400).JSON(fiber.Map{
			"message": "user email id is not verified",
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "user authentication failed",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	appSec := config.GetAppSecret()
	token, err := claims.SignedString([]byte(appSec))
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
		"message":  "login successful!",
		"username": user.Username,
		"id":       user.Id,
		"role":     user.Role,
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
	return c.Status(200).JSON(fiber.Map{
		"message": "logout suuccess!",
	})

}

func GetUser(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	appSec := config.GetAppSecret()

	if cookie != "" {
		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(appSec), nil
		})
		//fmt.Printf("%v", err)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message": "you are not authorized to perform this task",
			})
		}

		claims := token.Claims.(*jwt.StandardClaims)
		var user models.User

		database.DB.Where("id=?", claims.Issuer).First(&user)

		return c.Status(200).JSON(fiber.Map{
			"id":             user.Id,
			"fullname":       user.Username,
			"role":           user.Role,
			"email":          user.Email,
			"mobile":         user.Mobile,
			"email_verified": user.EmailVerified,
		})
	}

	return c.Status(400).JSON(fiber.Map{
		"data":    "",
		"message": "not received access token",
	})
}

func EmailVerify(c *fiber.Ctx) error {
	code := c.Params("code")

	var user models.User

	database.DB.Where("code", code).First(&user)
	if user.Id != 0 {
		user.EmailVerified = true
		database.DB.Save(&user)
		if database.DB.Error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "email is verification failed",
			})
		} else {
			return c.Status(200).JSON(fiber.Map{
				"message": "email is verified",
			})
		}
	} else {
		return c.Status(404).JSON(fiber.Map{
			"message": "user is not found",
		})
	}

}

func GetAllUsers(c *fiber.Ctx) error {
	var user []models.User
	database.DB.Order("created_at desc").Find(&user)

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}

func GetUserId(c *fiber.Ctx) (int, error) {
	cookie := c.Cookies("jwt")
	appSecret := config.GetAppSecret()
	if cookie != "" {
		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(appSecret), nil
		})
		fmt.Printf("%v", err)
		if err != nil {
			return 0, fmt.Errorf("you are not authorized %v", err)
		}

		claims := token.Claims.(*jwt.StandardClaims)

		var user models.User
		database.DB.Where("id=?", claims.Issuer).First(&user)

		return int(user.Id), nil
	}

	return 0, fmt.Errorf("user is not authorized")
}
