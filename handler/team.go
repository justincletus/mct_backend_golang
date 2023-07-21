package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/cms/database"
	"github.com/justincletus/cms/models"
)

func CreateTeam(c *fiber.Ctx) error {
	data := make(map[string]string)
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "bad request please check the payload",
		})
	}

	var user models.User

	database.DB.Where("role=?", data["role"]).First(&user)
	if user.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "user account not found",
		})
	}

	team := models.TeamMem{
		Title:           data["title"],
		SubContractor:   data["sub_contractor"],
		ContractorEmail: data["contractor_email"],
		ClientEmail:     data["client_email"],
		Members:         data["members"],
		UserId:          user.Id,
		User:            user,
	}

	database.DB.Create(&team)

	return c.Status(201).JSON(fiber.Map{
		"message": team,
	})

}

func GetTeams(c *fiber.Ctx) error {

	var teamMem []models.TeamMem
	//var user models.User

	database.DB.Order("created_at desc").Find(&teamMem)

	return c.Status(200).JSON(fiber.Map{
		"data": teamMem,
	})

}

func DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	var team models.TeamMem

	database.DB.Where("id", id).Find(&team)
	if team.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "team not found",
		})
	}

	database.DB.Unscoped().Delete(&team, id)

	return c.Status(200).JSON(fiber.Map{
		"message": "team is deleted",
	})

}
