package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"strconv"

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
	LogAction(c, "Labor", "Crear", "created labor id="+strconv.FormatUint(uint64(l.ID), 10))
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
	LogAction(c, "Labor", "Listar", "list labores")
	return c.JSON(labores)
}

func GetLabor(c *fiber.Ctx) error {
	id := c.Params("id")
	var l models.LaborAgronomica
	if err := database.DB.First(&l, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	LogAction(c, "Labor", "Obtener", "get labor id="+id)
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
	LogAction(c, "Labor", "Actualizar", "updated labor id="+id)
	return c.JSON(l)
}

func DeleteLabor(c *fiber.Ctx) error {
	id := c.Params("id")
	var l models.LaborAgronomica
	if err := database.DB.First(&l, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	database.DB.Delete(&l)
	LogAction(c, "Labor", "Eliminar", "deleted labor id="+id)
	return c.SendStatus(fiber.StatusNoContent)
}
