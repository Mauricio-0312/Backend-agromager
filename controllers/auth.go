package controllers

import (
	"strings"

	"agroproject/backend/database"
	"agroproject/backend/models"
	"agroproject/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
	Role     string `json:"role,omitempty"`
}

// SignUp crea usuario
func SignUp(c *fiber.Ctx) error {
	var body authRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	email := strings.ToLower(strings.TrimSpace(body.Email))
	if email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email y password requeridos"})
	}

	// Hash
	hash, err := models.HashPassword(body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error al hashear password"})
	}

	user := models.User{
		Email:    email,
		Password: hash,
		Name:     body.Name,
		Role:     body.Role,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "usuario ya existe"})
	}
	return c.JSON(fiber.Map{"message": "usuario creado"})
}

// Login -> retorna token
func Login(c *fiber.Ctx) error {
	var body authRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	email := strings.ToLower(strings.TrimSpace(body.Email))

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "credenciales inválidas"})
	}

	if !models.CheckPassword(body.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "credenciales inválidas"})
	}

	token, err := utils.GenerateToken(user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo generar token"})
	}

	return c.JSON(fiber.Map{"token": token, "role": user.Role})
}

// Me devuelve info simple del user
func Me(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.Claims)
	return c.JSON(fiber.Map{
		"email": claims.Email,
		"role":  claims.Role,
	})
}

// AdminGetUsers lista usuarios (solo admin)
func AdminGetUsers(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Select("id, email, name, role, created_at").Find(&users)
	return c.JSON(users)
}
