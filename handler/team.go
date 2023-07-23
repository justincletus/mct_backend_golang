package handler

import (
	"strconv"
	"strings"

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

	//var tmpMembers []string

	var flag = false

	if team.Members == "" {
		team.Members = data["member"]
	} else {
		if data["member"] == team.Members {
			return fiber.NewError(400, "user already exists")
		} else {
			if strings.Contains(team.Members, ",") {
				members := strings.Split(team.Members, ",")
				for _, member := range members {
					if member == data["member"] {
						flag = true
						break
					}
				}

				if flag {
					return fiber.NewError(500, "user already in group")
				} else {
					team.Members += "," + data["member"]
				}

			} else {
				team.Members += "," + data["member"]
			}
		}

	}

	database.DB.Save(&team)

	return c.Status(200).JSON(fiber.Map{
		"data": team,
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
