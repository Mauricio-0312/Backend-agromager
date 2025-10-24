package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Status      string         `gorm:"default:'activo'" json:"status"` // activo, cerrado, pausado
}
