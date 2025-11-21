package models

import (
	"time"

	"gorm.io/gorm"
)

// Logger records system events for audit
type Logger struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    *uint          `json:"user_id"` // user who caused the event (nullable)
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Module    string         `gorm:"size:100" json:"module"`
	Event     string         `gorm:"size:100" json:"event"`
	Details   string         `gorm:"type:text" json:"details"`
}
