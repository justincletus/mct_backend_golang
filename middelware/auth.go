package middelware

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/cms/config"
)

var SECRET = []byte(config.GetAppSecret())

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("jwt")
		if tokenString == "" {
			return fiber.ErrUnauthorized
		}

		//fmt.Println(SECRET)

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return SECRET, nil
		})

		if err != nil {
			fmt.Println(err)
			if err == jwt.ErrSignatureInvalid {
				return fiber.ErrUnauthorized
			}

			return fiber.ErrBadRequest
		}

		if !token.Valid {
			return fiber.ErrUnauthorized
		}
		return c.Next()
	}

}
