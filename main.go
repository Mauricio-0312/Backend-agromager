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
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "supersecret_jwt_key")
	}
	db := database.Connect()
	if db == nil {
		log.Fatal("DB connection failed")
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api := app.Group("/api")
	// auth
	api.Post("/login", controllers.Login)
	api.Post("/signup", controllers.SignUp)

	// protected routes
	protected := api.Group("/", middlewares.RequireAuth())

	// users (admin)
	protected.Get("/users", controllers.ListUsers)                       
	protected.Get("/users/:id", controllers.GetUser)
	protected.Put("/users/:id", controllers.UpdateUser)
	protected.Get("/csv/users", controllers.ExportUsersCSV)
	protected.Delete("/users/:id", controllers.DeleteUser)
	protected.Patch("/users/:id/password", controllers.ChangePassword)
	protected.Get("/me", controllers.Me)

	// projects
	protected.Post("/projects", controllers.CreateProject)
	protected.Get("/projects", controllers.ListProjects)
	protected.Get("/csv/projects", controllers.ExportProjectsCSV)
	protected.Get("/projects/:id", controllers.GetProject)
	protected.Put("/projects/:id", controllers.UpdateProject)
	protected.Patch("/projects/:id/close", controllers.CloseProject)

	// Labores Agronómicas
	protected.Post("/labores", controllers.CreateLabor)
	protected.Get("/labores", controllers.ListLabores)
	protected.Get("/labores/:id", controllers.GetLabor)
	protected.Put("/labores/:id", controllers.UpdateLabor)
	protected.Delete("/labores/:id", controllers.DeleteLabor)

	// Equipos e Implementos
	protected.Post("/equipos", controllers.CreateEquipo)
	protected.Get("/equipos", controllers.ListEquipos)
	protected.Get("/equipos/:id", controllers.GetEquipo)
	protected.Put("/equipos/:id", controllers.UpdateEquipo)
	protected.Delete("/equipos/:id", controllers.DeleteEquipo)

	// Actividades Agrícolas
	protected.Post("/activities", controllers.CreateActividad)
	protected.Get("/activities", controllers.ListActividades)
	protected.Get("/activities/:id", controllers.GetActividad)
	protected.Put("/activities/:id", controllers.UpdateActividad)
	protected.Delete("/activities/:id", controllers.DeleteActividad)

	// Units of Measure
	protected.Post("/units", controllers.CreateUnit)
	protected.Get("/units", controllers.ListUnits)
	protected.Get("/units/:id", controllers.GetUnit)
	protected.Put("/units/:id", controllers.UpdateUnit)
	protected.Delete("/units/:id", controllers.DeleteUnit)

	// Logger endpoints (admin only for listing/deleting)
	protected.Post("/logs", controllers.CreateLog)
	protected.Get("/logs", controllers.ListLogs)
	protected.Get("/logs/:id", controllers.GetLog)
	protected.Delete("/logs/:id", controllers.DeleteLog)
	protected.Get("/logs/count", controllers.CountLogs)

	log.Println("Backend listening on :8080")
	log.Fatal(app.Listen(":8080"))
}
