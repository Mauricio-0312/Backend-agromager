package controllers

import (
	"fmt"
	"strings"

	"agroproject/backend/database"
	"agroproject/backend/models"
	"agroproject/backend/utils"

	"github.com/gofiber/fiber/v2"
)

type signupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Dni      string `json:"dni"`
}

func SignUp(c *fiber.Ctx) error {
	var body signupReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	email := strings.ToLower(strings.TrimSpace(body.Email))
	if email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email y password requeridos"})
	}
	hash, _ := models.HashPassword(body.Password)
	user := models.User{
		Email:    email,
		Password: hash,
		Name:     body.Name,
		Dni:      body.Dni,
		Role:     body.Role,
		Active:   true,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "usuario ya existe"})
	}
	// Log signup event and provide user ID so logger embeds it even if not authenticated
	LogAction(c, "Usuario", "Registro", fmt.Sprintf("created user id=%d email=%s role=%s", user.ID, user.Email, user.Role), user.ID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "usuario creado", "user": user})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var body loginReq
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
	if !user.Active {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "usuario inactivo"})
	}
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo generar token"})
	}
	// Log login event, pass user ID so it is embedded
	LogAction(c, "Usuario", "Iniciar Sesión", fmt.Sprintf("user_id=%d email=%s", user.ID, user.Email), user.ID)
	return c.JSON(fiber.Map{"token": token, "role": user.Role, "user": fiber.Map{"id": user.ID, "email": user.Email, "name": user.Name, "role": user.Role, "dni": user.Dni}})
}
