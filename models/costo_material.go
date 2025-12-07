package models

import (
	"time"

	"gorm.io/gorm"
)

type CostoMaterialInsumo struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Actividad     string         `gorm:"size:255" json:"actividad"`
	Accion        string         `gorm:"size:500" json:"accion"`
	Categoria     int            `json:"categoria"` // 1=Materiales,2=Insumos,3=Equipos
	Descripcion   string         `gorm:"size:500" json:"descripcion"`
	Cantidad      float64        `json:"cantidad"`
	UnidadID      *uint          `json:"unidad_id"`
	Unidad        *UnitOfMeasure `gorm:"foreignKey:UnidadID" json:"unidad,omitempty"`
	Costo         float64        `json:"costo"`
	ResponsableID *uint          `json:"responsable_id"`
	Responsable   *User          `gorm:"foreignKey:ResponsableID" json:"responsable,omitempty"`
	Monto         float64        `json:"monto"`
	PlanAccionID  uint           `json:"plan_accion_id"`
}
