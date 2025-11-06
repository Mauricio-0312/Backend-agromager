package models

import (
	"time"

	"gorm.io/gorm"
)

// ActividadAgricola represents an activity tied to a project
type ActividadAgricola struct {
	ID                 uint               `gorm:"primaryKey" json:"id"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	DeletedAt          gorm.DeletedAt     `gorm:"index" json:"-"`
	Fecha              *time.Time         `json:"fecha"`
	Actividad          string             `json:"actividad"`
	LaborAgronomicaID  *uint              `json:"labor_agronomica_id"`
	LaborAgronomica    *LaborAgronomica   `gorm:"foreignKey:LaborAgronomicaID" json:"labor_agronomica,omitempty"`
	Equipos            []EquipoImplemento `gorm:"many2many:actividad_equipos;" json:"equipos,omitempty"`
	EncargadoID        *uint              `json:"encargado_id"`
	Encargado          *User              `gorm:"foreignKey:EncargadoID" json:"encargado,omitempty"`
	RecursoHumano      int                `json:"recurso_humano"`
	Costo              float64            `json:"costo"`
	Observaciones      string             `json:"observaciones"`
	ProjectID          uint               `json:"project_id"`
	Project            Project            `gorm:"foreignKey:ProjectID" json:"-"`
}
