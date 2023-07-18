package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/cms/database"
	"github.com/justincletus/cms/models"
)

func CreateTeam(c *fiber.Ctx) error {
	data := make(map[string]string)

	err := c.BodyParser(&data)
	//fmt.Println(data)
	if err != nil {
		//fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "bad request please check the payload",
		})
	}

	var user models.User

	database.DB.Where("email=?", data["manager"]).Where("role=?", "manager").Find(&user)
	if user.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "manager user email is not matched",
		})
	}

	team := models.TeamMem{
		Title:   data["title"],
		UserId:  user.Id,
		User:    user,
		Members: data["t_user"],
	}

	database.DB.Create(&team)

	return c.Status(201).JSON(fiber.Map{
		"message": "team is created",
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
