package database

import (
	"log"

	"agroproject/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("agro.db"), &gorm.Config{})
	if err != nil {
		log.Println("Error conectando a sqlite:", err)
		return nil
	}
	// Migraciones
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Println("AutoMigrate error:", err)
	}

	DB = db
	seedAdmin()
	return db
}

// seedAdmin crea un admin si no existe (contraseña por defecto: admin123) — cambiar en producción
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
		}
		DB.Create(&admin)
	}
}
