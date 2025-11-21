package database

import (
	"log"
	"os"

	"agroproject/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "agro.db"
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Println("Error conectando a sqlite:", err)
		return nil
	}

	// Migraciones
	err = db.AutoMigrate(&models.User{}, &models.Project{}, &models.UserProject{}, &models.LaborAgronomica{}, &models.EquipoImplemento{}, &models.ActividadAgricola{}, &models.Logger{}, &models.UnitOfMeasure{})
	if err != nil {
		log.Println("AutoMigrate error:", err)
	}

	DB = db
	seedAdmin()
	seedUnits()
	return db
}

func seedAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		pass, _ := models.HashPassword("admin123")
		admin := models.User{
			Email:    "admin@agro.local",
			Password: pass,
			Role:     "admin",
			Name:     "Administrador",
			Active:   true,
		}
		DB.Create(&admin)
	}
}

func seedUnits() {
	var count int64
	DB.Model(&models.UnitOfMeasure{}).Count(&count)
	if count == 0 {
		units := []models.UnitOfMeasure{
			{Dimension: "8", Unit: "Pulgadas"},
			{Dimension: "6", Unit: "Pulgadas"},
		}
		for _, u := range units {
			DB.Create(&u)
		}
	}
}
