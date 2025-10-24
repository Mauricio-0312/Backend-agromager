package models

type UserProject struct {
	UserID    uint   `gorm:"primaryKey"`
	ProjectID uint   `gorm:"primaryKey"`
	Role      string `json:"role" gorm:"default:'participant'"`
}
