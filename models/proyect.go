package models

import (
	"gorm.io/gorm"
)

type Proyect struct{
	gorm.Model
	Descripcion string `gorm:"not null"`
	Cierre *time.Time 
	Habilitado bool `gorm:"type:bool;default:true"`
}
