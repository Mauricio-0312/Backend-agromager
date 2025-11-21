package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"strconv"

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
	LogAction(c, "Equipo", "Crear", "created equipo id="+strconv.FormatUint(uint64(e.ID), 10))
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
	LogAction(c, "Equipo", "Listar", "list equipos")
	return c.JSON(items)
}

func GetEquipo(c *fiber.Ctx) error {
	id := c.Params("id")
	var e models.EquipoImplemento
	if err := database.DB.First(&e, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	LogAction(c, "Equipo", "Obtener", "get equipo id="+id)
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
	LogAction(c, "Equipo", "Actualizar", "updated equipo id="+id)
	return c.JSON(e)
}

func DeleteEquipo(c *fiber.Ctx) error {
	id := c.Params("id")
	var e models.EquipoImplemento
	if err := database.DB.First(&e, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	database.DB.Delete(&e)
	LogAction(c, "Equipo", "Eliminar", "deleted equipo id="+id)
	return c.SendStatus(fiber.StatusNoContent)
}
