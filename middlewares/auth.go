package middlewares

import (
	"strings"

	"agroproject/backend/utils"

	"github.com/gofiber/fiber/v2"
)

// Protected middleware valida JWT y coloca claims en Locals
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authorization header missing"})
		}
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authorization header invalid"})
		}
		tokenStr := parts[1]
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token inválido"})
		}
		c.Locals("claims", claims)
		return c.Next()
	}
}

// AdminOnly middleware verifica que role == "admin"
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*utils.Claims)
		if claims.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "recurso solo para administradores"})
		}
		return c.Next()
	}
}
