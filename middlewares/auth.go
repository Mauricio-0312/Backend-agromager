package middlewares

import (
	"strings"

	"agroproject/backend/utils"

	"github.com/gofiber/fiber/v2"
)

// RequireAuth verifica JWT y pone claims en locals
func RequireAuth() fiber.Handler {
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
		c.Locals("user_id", claims.UserID)
		c.Locals("user_role", claims.Role)
		return c.Next()
	}
}

// RequireRole verifica rol mínimo
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		r := c.Locals("user_role")
		if r == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no role found"})
		}
		userRole := r.(string)
		if userRole != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "acceso prohibido: se requiere rol " + role})
		}
		return c.Next()
	}
}
