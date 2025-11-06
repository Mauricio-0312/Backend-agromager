package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"

	"github.com/gofiber/fiber/v2"
)

type laborReq struct {
	Descripcion string `json:"descripcion"`
}

func CreateLabor(c *fiber.Ctx) error {
	var body laborReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid"})
	}
	l := models.LaborAgronomica{Descripcion: body.Descripcion}
	if err := database.DB.Create(&l).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create"})
	}
	return c.JSON(l)
}

func ListLabores(c *fiber.Ctx) error {
	var labores []models.LaborAgronomica
	q := c.Query("q")
	db := database.DB.Model(&models.LaborAgronomica{})
	if q != "" {
		db = db.Where("descripcion LIKE ?", "%"+q+"%")
	}
	db.Find(&labores)
	return c.JSON(labores)
}

func GetLabor(c *fiber.Ctx) error {
	id := c.Params("id")
	var l models.LaborAgronomica
	if err := database.DB.First(&l, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(l)
}

func UpdateLabor(c *fiber.Ctx) error {
	id := c.Params("id")
	var body laborReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid"})
	}
	var l models.LaborAgronomica
	if err := database.DB.First(&l, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	if body.Descripcion != "" {
		l.Descripcion = body.Descripcion
	}
	database.DB.Save(&l)
	return c.JSON(l)
}

func DeleteLabor(c *fiber.Ctx) error {
	id := c.Params("id")
	var l models.LaborAgronomica
	if err := database.DB.First(&l, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	database.DB.Delete(&l)
	return c.SendStatus(fiber.StatusNoContent)
}
