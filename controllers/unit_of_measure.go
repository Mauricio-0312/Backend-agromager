package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Create unit
func CreateUnit(c *fiber.Ctx) error {
	var u models.UnitOfMeasure
	if err := c.BodyParser(&u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	if err := database.DB.Create(&u).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo crear unidad"})
	}
	// Log creation
	LogAction(c, "Unidad de medida", "Crear", fmt.Sprintf("id=%d unit=%s dimension=%s", u.ID, u.Unit, u.Dimension))
	return c.Status(fiber.StatusCreated).JSON(u)
}

// List units
func ListUnits(c *fiber.Ctx) error {
	var units []models.UnitOfMeasure
	database.DB.Find(&units)

	// Log listing
	LogAction(c, "Unidad de medida", "Listar", fmt.Sprintf("count=%d", len(units)))
	return c.JSON(units)
}

// Get unit
func GetUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	var u models.UnitOfMeasure
	if err := database.DB.First(&u, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "unidad no encontrada"})
	}

	LogAction(c, "Unidad de medida", "Obtener", fmt.Sprintf("id=%d unit=%s dimension=%s", u.ID, u.Unit, u.Dimension))
	return c.JSON(u)
}

// Update unit
func UpdateUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	var body models.UnitOfMeasure
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	var u models.UnitOfMeasure
	if err := database.DB.First(&u, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "unidad no encontrada"})
	}
	if body.Dimension != "" {
		u.Dimension = body.Dimension
	}
	if body.Unit != "" {
		u.Unit = body.Unit
	}
	database.DB.Save(&u)
	// Log update
	LogAction(c, "Unidad de medida", "Actualizar", fmt.Sprintf("id=%d unit=%s dimension=%s", u.ID, u.Unit, u.Dimension))
	return c.JSON(u)
}

// Delete unit
func DeleteUnit(c *fiber.Ctx) error {
	id := c.Params("id")
	// fetch unit to include info in log
	var u models.UnitOfMeasure
	if err := database.DB.First(&u, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "unidad no encontrada"})
	}
	if err := database.DB.Delete(&models.UnitOfMeasure{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo eliminar unidad"})
	}
	// Log deletion
	LogAction(c, "Unidad de medida", "Eliminar", fmt.Sprintf("id=%d unit=%s dimension=%s", u.ID, u.Unit, u.Dimension))
	return c.JSON(fiber.Map{"message": "unidad eliminada"})
}
