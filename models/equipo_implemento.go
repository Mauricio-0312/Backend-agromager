package models

import (
	"time"

	"gorm.io/gorm"
)

// EquipoImplemento represents equipment or implement
type EquipoImplemento struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Descripcion string       `gorm:"not null" json:"descripcion"`
}
