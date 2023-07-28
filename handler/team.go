package handler

import (
	"strconv"

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

	// var members []string
	// members = append(members, data["members"])

	team := models.TeamMem{
		Title:           data["title"],
		SubContractor:   data["sub_contractor"],
		ContractorEmail: data["contractor_email"],
		ClientEmail:     data["client_email"],
		Members:         data["client_inspector"],
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
	var members []models.Member

	database.DB.Order("created_at desc").Find(&teamMem)
	database.DB.Order("created_at desc").Find(&members)

	return c.Status(200).JSON(fiber.Map{
		"data":    teamMem,
		"members": members,
	})

}

func GetTeamById(c *fiber.Ctx) error {
	id := c.Params("id")

	var team models.TeamMem
	database.DB.Where("id=?", id).First(&team)
	if team.Id == 0 {
		return fiber.NewError(404, "team by id not found")
	}

	var users []models.User

	database.DB.Order("id desc").Where("role=?", "user").Find(&users)

	return c.Status(200).JSON(fiber.Map{
		"data":  team,
		"users": users,
	})

}

func UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	tID, _ := strconv.ParseInt(id, 10, 32)

	var team models.TeamMem
	database.DB.Where("id=?", int(tID)).First(&team)
	if team.Id == 0 {
		return fiber.NewError(404, "team with id is not found")
	}

	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return fiber.NewError(500, "error in post data")
	}

	var member models.Member
	member.Email = data["member"]
	member.TeamId = team.Id

	database.DB.Create(&member)

	if member.Id == 0 {
		return fiber.NewError(500, "add member to the team failed")
	}

	var members []models.Member
	database.DB.Order("created_at desc").Find(&members)

	return c.Status(200).JSON(fiber.Map{
		"data":    team,
		"members": members,
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
