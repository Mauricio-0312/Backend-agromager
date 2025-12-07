package models

import (
	"time"

	"gorm.io/gorm"
)

type PlanAccion struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Actividad      string         `gorm:"size:255" json:"actividad"`
	Accion         string         `gorm:"size:500" json:"accion"`
	FechaInicio    *time.Time     `json:"fecha_inicio"`
	FechaCierre    *time.Time     `json:"fecha_cierre"`
	CantidadHoras  int            `json:"cantidad_horas"`
	ResponsableID  *uint          `json:"responsable_id"`
	Responsable    *User          `gorm:"foreignKey:ResponsableID" json:"responsable,omitempty"`
	Monto          float64        `json:"monto"`
	ProyectoID     uint           `json:"proyecto_id"`
}
