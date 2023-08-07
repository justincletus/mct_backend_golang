package handler

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/mct_backend_golang/database"
	"github.com/justincletus/mct_backend_golang/models"
)

func CreateJob(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)
	if err != nil {
		fmt.Errorf("%s", err.Error())
		return c.Status(400).JSON(fiber.Map{
			"message": "not able received data ",
		})
	}

	// fmt.Println(data)

	if data["name"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "project should not empty",
		})
	}

	u_id, err := GetUserId(c)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	job := models.Job{
		Name:   data["name"],
		JobId:  data["job_id"],
		UserId: u_id,
	}

	txtID := database.DB.Create(&job)
	id := txtID.RowsAffected
	if id == 0 {
		message := fmt.Sprintf("%v", txtID.Error)
		code := strings.Contains(message, "1062")

		if code {
			return c.Status(400).JSON(fiber.Map{
				"message": "duplicate entry of job name",
			})
		} else {
			return c.Status(400).JSON(fiber.Map{
				"message": "data not saved please check the data and resubmit",
			})
		}

	}

	return c.Status(200).JSON(fiber.Map{
		"id":     job.Id,
		"job_id": job.JobId,
	})

}

func GetJobs(c *fiber.Ctx) error {
	var jobs []models.Job

	database.DB.Order("created_at desc").Find(&jobs)

	return c.Status(200).JSON(fiber.Map{
		"data": jobs,
	})
}

func GetProjects(c *fiber.Ctx) error {
	var projects []models.Job

	database.DB.Order("created_at desc").Find(&projects)

	return c.Status(200).JSON(fiber.Map{
		"data": projects,
	})
}

func DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	//fmt.Println(id)
	var project models.Job
	database.DB.Where("id=?", id).First(&project)
	fmt.Println(project)

	if project.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "project not found",
		})

	}
	database.DB.Unscoped().Delete(&project, id)
	return c.Status(204).JSON(fiber.Map{})
}

func UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Job

	var data map[string]string

	database.DB.Where("id=?", id).First(&project)
	if project.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "post data not able process",
		})

	}

	if data["p_name"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "project name should not empty",
		})
	}

	project.Name = data["p_name"]

	database.DB.Save(&project)

	return c.Status(200).JSON(fiber.Map{
		"message": "project details updated",
	})

}

func GetProjectByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var project models.Job

	database.DB.Where("id=?", id).First(&project)
	if project.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "project not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": project,
	})
}
