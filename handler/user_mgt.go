package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/cms/database"
	"github.com/justincletus/cms/models"
)

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	//fmt.Println(id)

	var user models.User

	database.DB.Where("id", id).First(&user)
	if user.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "user is not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}

func UpdateUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "bad request, please check the post data",
		})
	}

	var user models.User
	var manager models.Manager

	database.DB.Where("id", id).First(&user)
	if user.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if data["fullname"] != "" {
		user.Fullname = data["fullname"]
	}
	if data["email"] != "" {
		user.Email = data["email"]
	}
	if data["mobile"] != "" {
		user.Mobile = data["mobile"]
	}
	if data["role"] != "" {
		user.Role = data["role"]
		if data["role"] == "manager" {
			database.DB.Where("user_id", user.Id).First(&manager)
			fmt.Println(manager)
			if manager.Id == 0 {
				manager.UserId = user.Id
				database.DB.Create(&manager)
			}

		} else if data["role"] == "user" {
			database.DB.Where("user_id", user.Id).First(&manager)
			if manager.Id != 0 {
				database.DB.Unscoped().Delete(&manager, manager.Id)
			}
		}

	}

	if data["email_verifield"] != "" {
		if data["email_verified"] == "true" {
			user.EmailVerified = true
		} else {
			user.EmailVerified = false
		}
	}

	database.DB.Save(&user)

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})

}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	database.DB.Where("id", id).First(&user)
	if user.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}
	database.DB.Unscoped().Delete(&user, id)

	return c.Status(200).JSON(fiber.Map{
		"message": "user is deleted",
	})
}
