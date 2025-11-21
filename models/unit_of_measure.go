package models

import (
	"time"

	"gorm.io/gorm"
)

type UnitOfMeasure struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Dimension string         `gorm:"size:100;not null" json:"dimension"`
	Unit      string         `gorm:"size:100;not null" json:"unit"`
}
