package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateUuid() string {
	uuid := uuid.New()
	return uuid.String()
}

func GetRemoteHostAddress(c *fiber.Ctx) string {
	//ipAddress := c.Request().Host()

	host := c.Context().RemoteAddr().String()
	//host1 := host.String()
	var ipAddress string
	var ipAddressPort []string
	status := strings.Contains(host, ":")
	if status {
		ipAddressPort = strings.Split(host, ":")
		ipAddress = ipAddressPort[0] + ":3000"
		return ipAddress
	} else {
		return host
	}

}
