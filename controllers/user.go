package controllers

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"time"

	"agroproject/backend/database"
	"agroproject/backend/models"
	"agroproject/backend/utils"

	"github.com/gofiber/fiber/v2"
	
)

// List users (admin)
func ListUsers(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Select("id, email, name, role, active, created_at").Find(&users)
	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "usuario no encontrado"})
	}
	return c.JSON(user)
}

type updateUserReq struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Active *bool `json:"active"`
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var body updateUserReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "usuario no encontrado"})
	}
	if body.Name != "" {
		user.Name = body.Name
	}
	if body.Role != "" {
		user.Role = body.Role
	}
	if body.Active != nil {
		user.Active = *body.Active
	}
	database.DB.Save(&user)
	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo eliminar usuario"})
	}
	return c.JSON(fiber.Map{"message": "usuario eliminado"})
}

func ExportUsersCSV(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Find(&users)

	// create csv
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	// header
	writer.Write([]string{"id", "email", "name", "role", "active", "created_at"})
	for _, u := range users {
		created := u.CreatedAt.Format(time.RFC3339)
		writer.Write([]string{
			strconv.FormatUint(uint64(u.ID), 10),
			u.Email,
			u.Name,
			u.Role,
			strconv.FormatBool(u.Active),
			created,
		})
	}
	writer.Flush()

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=users.csv")
	return c.SendStream(&buf)
}

// Change password
type changePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	var req changePasswordReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "usuario no encontrado"})
	}
	if !models.CheckPassword(req.OldPassword, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "contraseña actual incorrecta"})
	}
	hash, _ := models.HashPassword(req.NewPassword)
	user.Password = hash
	database.DB.Save(&user)
	return c.JSON(fiber.Map{"message": "password updated"})
}

func Me(c *fiber.Ctx) error {
	claimsVal := c.Locals("claims")
	if claimsVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	claims, ok := claimsVal.(*utils.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(fiber.Map{
		"email": claims.Email,
		"role":  claims.Role,
	})
}