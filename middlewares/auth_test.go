package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"agroproject/backend/utils"
	"github.com/gofiber/fiber/v2"
)

func TestProtectedMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(Protected())
	app.Get("/test", func(c *fiber.Ctx) error {
		claims := c.Locals("claims")
		if claims == nil {
			return c.Status(401).SendString("No claims")
		}
		return c.SendString("OK")
	})

	// Token válido
	token, _ := utils.GenerateToken("test@example.com", "user")
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error en petición protegida: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Middleware falló con token válido, status: %d", resp.StatusCode)
	}
}
