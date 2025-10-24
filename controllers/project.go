package controllers

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"time"

	"agroproject/backend/database"
	"agroproject/backend/models"

	"github.com/gofiber/fiber/v2"
)

type projectReq struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      string     `json:"status"`
}

func CreateProject(c *fiber.Ctx) error {
	var body projectReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid"})
	}
	p := models.Project{
		Name:        body.Name,
		Description: body.Description,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
		Status:      "activo",
	}
	if body.Status != "" {
		p.Status = body.Status
	}
	if err := database.DB.Create(&p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create"})
	}
	return c.JSON(p)
}

func ListProjects(c *fiber.Ctx) error {
	var projects []models.Project
	q := c.Query("q")
	db := database.DB
	if q != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+q+"%", "%"+q+"%")
	}
	db.Find(&projects)
	return c.JSON(projects)
}

func GetProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var p models.Project
	if err := database.DB.First(&p, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
	}
	return c.JSON(p)
}

func UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var body projectReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid"})
	}
	var p models.Project
	if err := database.DB.First(&p, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
	}
	if body.Name != "" {
		p.Name = body.Name
	}
	if body.Description != "" {
		p.Description = body.Description
	}
	if body.StartDate != nil {
		p.StartDate = body.StartDate
	}
	if body.EndDate != nil {
		p.EndDate = body.EndDate
	}
	if body.Status != "" {
		p.Status = body.Status
	}
	database.DB.Save(&p)
	return c.JSON(p)
}

func CloseProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var p models.Project
	if err := database.DB.First(&p, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
	}
	p.Status = "cerrado"
	now := time.Now()
	p.EndDate = &now
	database.DB.Save(&p)
	return c.JSON(p)
}

func ExportProjectsCSV(c *fiber.Ctx) error {
	var projects []models.Project
	database.DB.Find(&projects)
	// build csv buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"id", "name", "description", "start_date", "end_date", "status", "created_at"})
	for _, p := range projects {
		var sd, ed string
		if p.StartDate != nil {
			sd = p.StartDate.Format(time.RFC3339)
		}
		if p.EndDate != nil {
			ed = p.EndDate.Format(time.RFC3339)
		}
		writer.Write([]string{
			strconv.FormatUint(uint64(p.ID), 10),
			p.Name,
			p.Description,
			sd,
			ed,
			p.Status,
			p.CreatedAt.Format(time.RFC3339),
		})
	}
	writer.Flush()
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=projects.csv")
	return c.SendStream(&buf)
}
