package handler

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/cms/config"
	"github.com/justincletus/cms/database"
	"github.com/justincletus/cms/models"
	"github.com/k3a/html2text"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
	"gorm.io/gorm"
)

var SECRET = config.SECRET

var temp *template.Template

func init() {
	temp = template.Must(template.ParseGlob("templates/*.html"))
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	var user models.User

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "user registration failed",
		})
	}

	// fmt.Println(data["email"])

	// err = sendEmail(data["email"])
	// if err != nil {
	// 	return c.Status(400).JSON(fiber.Map{
	// 		"message": "sending email is failed",
	// 	})
	// }

	// return c.Status(200).JSON(fiber.Map{
	// 	"message": "email is send successfully",
	// })

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
	if data["role"] == "" {
		data["role"] = "user"
	}

	user.Role = data["role"]

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

	return c.Status(201).JSON(fiber.Map{
		"data": user,
	})

}

func getUsername(name string) string {
	name = strings.ToLower(name)
	if strings.Contains(name, " ") {
		str := strings.Split(name, " ")
		name = strings.Join(str, "")
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

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
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

	if cookie != "" {
		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET), nil
		})
		fmt.Printf("%v", err)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message": "you are not authorized to perform this task",
			})
		}

		claims := token.Claims.(*jwt.StandardClaims)
		var user models.User

		fmt.Println(claims.Issuer)

		database.DB.Where("id=?", claims.Issuer).First(&user)
		return c.Status(200).JSON(fiber.Map{
			"id":       user.ID,
			"fullname": user.Fullname,
		})
	}

	return c.Status(400).JSON(fiber.Map{
		"message": "not received access token",
	})
}

func sendEmail(email string) error {
	from := "info@cccabs-service.in"
	password := "#03Jgabi"

	smtpHost := "smtp.hostinger.com"
	smtpPort := 587

	var body bytes.Buffer

	data := struct {
		Name    string
		Message string
	}{
		Name:    "Hello",
		Message: "World",
	}

	if err := temp.ExecuteTemplate(&body, "template.html", &data); err != nil {
		fmt.Println(err)
		return err
	}

	// return nil
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Demo email")
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil

}
