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
	UserIDs     []uint     `json:"user_ids"`
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
	// associate users if provided
	if len(body.UserIDs) > 0 {
		var users []models.User
		database.DB.Find(&users, body.UserIDs)
		if len(users) > 0 {
			database.DB.Model(&p).Association("Users").Replace(users)
		}
	}
	return c.JSON(p)
}

func ListProjects(c *fiber.Ctx) error {
	var projects []models.Project
	q := c.Query("q")
	db := database.DB.Model(&models.Project{})
	if q != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+q+"%", "%"+q+"%")
	}

	// Restrict results for non-admin users: only projects where the user is a participant
	roleLoc := c.Locals("user_role")
	userIDLoc := c.Locals("user_id")
	if roleStr, ok := roleLoc.(string); ok && roleStr != "admin" {
		// coerce user id to uint (support several possible underlying types)
		var uid uint
		switch v := userIDLoc.(type) {
		case uint:
			uid = v
		case int:
			uid = uint(v)
		case int64:
			uid = uint(v)
		case float64:
			uid = uint(v)
		default:
			uid = 0
		}
		// join the linking table to filter projects where this user is linked
		db = db.Joins("JOIN user_projects up ON up.project_id = projects.id").Where("up.user_id = ?", uid)
	}

	db.Preload("Users").Find(&projects)
	return c.JSON(projects)
}

func GetProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var p models.Project
	if err := database.DB.Preload("Users").First(&p, id).Error; err != nil {
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
	// update associations if provided
	if body.UserIDs != nil {
		var users []models.User
		database.DB.Find(&users, body.UserIDs)
		database.DB.Model(&p).Association("Users").Replace(users)
	}
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

	// Build base DB query
	db := database.DB.Model(&models.Project{})

	// Restrict export for non-admin users: only include projects where the user is a participant
	roleLoc := c.Locals("user_role")
	userIDLoc := c.Locals("user_id")
	if roleStr, ok := roleLoc.(string); ok && roleStr != "admin" {
		var uid uint
		switch v := userIDLoc.(type) {
		case uint:
			uid = v
		case int:
			uid = uint(v)
		case int64:
			uid = uint(v)
		case float64:
			uid = uint(v)
		default:
			uid = 0
		}
		// join linking table to filter by participant
		db = db.Joins("JOIN user_projects up ON up.project_id = projects.id").Where("up.user_id = ?", uid)
	}

	db.Preload("Users").Find(&projects)
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
