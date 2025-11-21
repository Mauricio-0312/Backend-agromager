package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CreateLog adds a log entry
func CreateLog(c *fiber.Ctx) error {
	var logm models.Logger
	if err := c.BodyParser(&logm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	if err := database.DB.Create(&logm).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo crear log"})
	}
	return c.JSON(logm)
}

// List logs with optional query
func ListLogs(c *fiber.Ctx) error {
	var logs []models.Logger
	q := c.Query("q")
	db := database.DB.Model(&models.Logger{})
	if q != "" {
		db = db.Where("module LIKE ? OR event LIKE ? OR details LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}
	// preload user and order by creation time descending
	db.Preload("User").Order("created_at desc").Find(&logs)
	return c.JSON(logs)
}

// Get single log
func GetLog(c *fiber.Ctx) error {
	id := c.Params("id")
	var l models.Logger
	if err := database.DB.Preload("User").First(&l, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "log no encontrado"})
	}
	return c.JSON(l)
}

// Delete log
func DeleteLog(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.Logger{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo eliminar log"})
	}
	return c.JSON(fiber.Map{"message": "log eliminado"})
}

// Count logs
func CountLogs(c *fiber.Ctx) error {
	var count int64
	database.DB.Model(&models.Logger{}).Count(&count)
	return c.JSON(fiber.Map{"count": strconv.FormatInt(count, 10)})
}
