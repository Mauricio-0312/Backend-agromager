package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type actividadReq struct {
	Fecha             *time.Time `json:"fecha"`
	Actividad         string     `json:"actividad"`
	LaborAgronomicaID *uint      `json:"labor_agronomica_id"`
	EquiposIDs        []uint     `json:"equipos_ids"`
	EncargadoID       *uint      `json:"encargado_id"`
	RecursoHumano     int        `json:"recurso_humano"`
	Costo             float64    `json:"costo"`
	Observaciones     string     `json:"observaciones"`
	ProjectID         uint       `json:"project_id"`
}

func CreateActividad(c *fiber.Ctx) error {
	var body actividadReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid"})
	}
	a := models.ActividadAgricola{
		Fecha:             body.Fecha,
		Actividad:         body.Actividad,
		LaborAgronomicaID: body.LaborAgronomicaID,
		EncargadoID:       body.EncargadoID,
		RecursoHumano:     body.RecursoHumano,
		Costo:             body.Costo,
		Observaciones:     body.Observaciones,
		ProjectID:         body.ProjectID,
	}
	if err := database.DB.Create(&a).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create"})
	}
	// associate equipos if provided
	if len(body.EquiposIDs) > 0 {
		var equipos []models.EquipoImplemento
		database.DB.Find(&equipos, body.EquiposIDs)
		if len(equipos) > 0 {
			database.DB.Model(&a).Association("Equipos").Replace(equipos)
		}
	}
	// preload relations
	database.DB.Preload("Equipos").Preload("LaborAgronomica").Preload("Encargado").First(&a, a.ID)
	LogAction(c, "Actividad", "Crear", "created activity id="+strconv.FormatUint(uint64(a.ID), 10))
	return c.JSON(a)
}

func ListActividades(c *fiber.Ctx) error {
	var acts []models.ActividadAgricola
	q := c.Query("q")
	db := database.DB.Model(&models.ActividadAgricola{})
	if q != "" {
		db = db.Where("actividad LIKE ?", "%"+q+"%")
	}
	// optional project filter
	if pid := c.Query("project_id"); pid != "" {
		db = db.Where("project_id = ?", pid)
	}
	db.Preload("Equipos").Preload("LaborAgronomica").Preload("Encargado").Find(&acts)
	LogAction(c, "Actividad", "Listar", "list activities")
	return c.JSON(acts)
}

func GetActividad(c *fiber.Ctx) error {
	id := c.Params("id")
	var a models.ActividadAgricola
	if err := database.DB.Preload("Equipos").Preload("LaborAgronomica").Preload("Encargado").First(&a, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	LogAction(c, "Actividad", "Obtener", "get activity id="+id)
	return c.JSON(a)
}

func UpdateActividad(c *fiber.Ctx) error {
	id := c.Params("id")
	var body actividadReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid"})
	}
	var a models.ActividadAgricola
	if err := database.DB.First(&a, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	if body.Fecha != nil {
		a.Fecha = body.Fecha
	}
	if body.Actividad != "" {
		a.Actividad = body.Actividad
	}
	if body.LaborAgronomicaID != nil {
		a.LaborAgronomicaID = body.LaborAgronomicaID
	}
	if body.EncargadoID != nil {
		a.EncargadoID = body.EncargadoID
	}
	if body.RecursoHumano != 0 {
		a.RecursoHumano = body.RecursoHumano
	}
	if body.Costo != 0 {
		a.Costo = body.Costo
	}
	if body.Observaciones != "" {
		a.Observaciones = body.Observaciones
	}
	database.DB.Save(&a)
	if body.EquiposIDs != nil {
		var equipos []models.EquipoImplemento
		database.DB.Find(&equipos, body.EquiposIDs)
		database.DB.Model(&a).Association("Equipos").Replace(equipos)
	}
	database.DB.Preload("Equipos").Preload("LaborAgronomica").Preload("Encargado").First(&a, a.ID)
	LogAction(c, "Actividad", "Actualizar", "updated activity id="+id)
	return c.JSON(a)
}

func DeleteActividad(c *fiber.Ctx) error {
	id := c.Params("id")
	var a models.ActividadAgricola
	if err := database.DB.First(&a, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	database.DB.Delete(&a)
	LogAction(c, "Actividad", "Eliminar", "deleted activity id="+id)
	return c.SendStatus(fiber.StatusNoContent)
}
