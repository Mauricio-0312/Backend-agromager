package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"

	"github.com/gofiber/fiber/v2"
)

type equipoReq struct {
	Descripcion string `json:"descripcion"`
}

func CreateEquipo(c *fiber.Ctx) error {
	var body equipoReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid"})
	}
	e := models.EquipoImplemento{Descripcion: body.Descripcion}
	if err := database.DB.Create(&e).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create"})
	}
	return c.JSON(e)
}

func ListEquipos(c *fiber.Ctx) error {
	var items []models.EquipoImplemento
	q := c.Query("q")
	db := database.DB.Model(&models.EquipoImplemento{})
	if q != "" {
		db = db.Where("descripcion LIKE ?", "%"+q+"%")
	}
	db.Find(&items)
	return c.JSON(items)
}

func GetEquipo(c *fiber.Ctx) error {
	id := c.Params("id")
	var e models.EquipoImplemento
	if err := database.DB.First(&e, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(e)
}

func UpdateEquipo(c *fiber.Ctx) error {
	id := c.Params("id")
	var body equipoReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid"})
	}
	var e models.EquipoImplemento
	if err := database.DB.First(&e, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	if body.Descripcion != "" {
		e.Descripcion = body.Descripcion
	}
	database.DB.Save(&e)
	return c.JSON(e)
}

func DeleteEquipo(c *fiber.Ctx) error {
	id := c.Params("id")
	var e models.EquipoImplemento
	if err := database.DB.First(&e, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	database.DB.Delete(&e)
	return c.SendStatus(fiber.StatusNoContent)
}
