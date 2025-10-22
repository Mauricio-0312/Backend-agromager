package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"agroproject/backend/database"
)

func setupTestDB() {
	database.Connect()
}

func TestSignUpHandler(t *testing.T) {
	setupTestDB()
	app := fiber.New()
	app.Post("/signup", SignUp)
	payload := `{"email":"test@agro.local","password":"123456","name":"Test"}`
	req := httptest.NewRequest("POST", "/signup", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error en petición signup: %v", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("SignUp falló, status: %d", resp.StatusCode)
	}
}

func TestLoginHandler(t *testing.T) {
	setupTestDB()
	app := fiber.New()
	app.Post("/login", Login)
	payload := `{"email":"test@agro.local","password":"123456"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error en petición login: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Login falló, status: %d", resp.StatusCode)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if _, ok := result["token"]; !ok {
		t.Error("No se recibió token en respuesta de login")
	}
}
