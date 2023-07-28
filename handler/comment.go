package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/cms/database"
	"github.com/justincletus/cms/models"
)

func GetComments(c *fiber.Ctx) error {
	var comments []models.Comment

	database.DB.Order("id desc").Find(&comments)

	return c.Status(200).JSON(fiber.Map{
		"data": comments,
	})
}

func DeleteComment(c *fiber.Ctx) error {
	id := c.Params("id")

	var comment models.Comment

	database.DB.Where("id=?", id).First(&comment)

	if comment.Id == 0 {
		return fiber.NewError(404, "comment not found")
	}

	database.DB.Unscoped().Delete(comment, id)

	return c.Status(204).JSON(fiber.Map{
		"message": "comment deleted",
	})
}
