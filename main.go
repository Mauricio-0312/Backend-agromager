package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"agroproject/backend/database"
	"agroproject/backend/controllers"
	"agroproject/backend/middlewares"
)

func main() {
	// Load JWT secret from env or default
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "supersecret_jwt_key_change_me")
	}

	// Init DB
	db := database.Connect()
	if db == nil {
		log.Fatal("No se pudo conectar a la base de datos")
	}

	// Fiber app
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // Vite default dev URL - ajusta si es necesario
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	api := app.Group("/api")
	api.Post("/signup", controllers.SignUp)
	api.Post("/login", controllers.Login)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middlewares.Protected()) // middleware que valida JWT
	protected.Get("/me", controllers.Me)

	// Admin-only
	admin := protected.Group("/")
	admin.Use(middlewares.AdminOnly())
	admin.Get("/admin/users", controllers.AdminGetUsers)

	// Start
	log.Println("Backend escuchando en :8080")
	log.Fatal(app.Listen(":8080"))
}
