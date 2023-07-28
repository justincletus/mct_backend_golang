package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/cms/database"
	"github.com/justincletus/cms/models"
)

func GetClientReports(c *fiber.Ctx) error {
	var c_reports []models.ClientReport

	database.DB.Order("created_at desc").Find(&c_reports)

	return c.Status(200).JSON(fiber.Map{
		"data": c_reports,
	})
}

func DeleteClientReport(c *fiber.Ctx) error {
	id := c.Params("id")
	var c_report models.ClientReport

	database.DB.Where("id=?", id).First(&c_report)

	if c_report.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "client report not found",
		})
	}

	database.DB.Unscoped().Delete(c_report, id)

	return c.Status(204).JSON(fiber.Map{
		"message": "report deleted!",
	})
}
