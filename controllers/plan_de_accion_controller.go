package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type planReq struct {
	Actividad     string  `json:"actividad"`
	Accion        string  `json:"accion"`
	FechaInicio   *string `json:"fecha_inicio"`
	FechaCierre   *string `json:"fecha_cierre"`
	CantidadHoras *int    `json:"cantidad_horas"`
	ResponsableID *uint   `json:"responsable_id"`
	Monto         *float64 `json:"monto"`
}

// CreatePlanAccion creates a plan under a project
func CreatePlanAccion(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var req planReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}

	// debug log incoming data
	fmt.Println("CreatePlanAccion request - projectID:", projectID, "payload:", req)

	// build model
	var body models.PlanAccion
	body.Actividad = req.Actividad
	body.Accion = req.Accion
	if req.CantidadHoras != nil { body.CantidadHoras = *req.CantidadHoras }
	if req.ResponsableID != nil { body.ResponsableID = req.ResponsableID }
	// monto is read-only and will be calculated from costos; initialize to 0
	body.Monto = 0

	// parse dates if present
	if req.FechaInicio != nil && *req.FechaInicio != "" {
		if t, err := time.Parse(time.RFC3339, *req.FechaInicio); err == nil {
			body.FechaInicio = &t
		} else if t2, err2 := time.Parse("2006-01-02", *req.FechaInicio); err2 == nil {
			body.FechaInicio = &t2
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fecha_inicio inválida"})
		}
	}
	if req.FechaCierre != nil && *req.FechaCierre != "" {
		if t, err := time.Parse(time.RFC3339, *req.FechaCierre); err == nil {
			body.FechaCierre = &t
		} else if t2, err2 := time.Parse("2006-01-02", *req.FechaCierre); err2 == nil {
			body.FechaCierre = &t2
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fecha_cierre inválida"})
		}
	}

	// validate dates
	if body.FechaInicio != nil && body.FechaCierre != nil {
		if body.FechaCierre.Before(*body.FechaInicio) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fecha_cierre debe ser >= fecha_inicio"})
		}
	}
	// cant horas >= 0
	if body.CantidadHoras < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cantidad_horas inválida"})
	}

	// attach project id
	var pid uint
	if _, err := fmt.Sscanf(projectID, "%d", &pid); err != nil {
		fmt.Println("CreatePlanAccion: failed to parse projectID:", projectID, "error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "project id inválido"})
	}
	fmt.Println("CreatePlanAccion: parsed pid:", pid)
	body.ProyectoID = pid

	// validate project exists and active
	var p models.Project
	if err := database.DB.First(&p, pid).Error; err != nil {
		fmt.Println("CreatePlanAccion: project not found pid:", pid, "db error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "proyecto no encontrado"})
	}
	if p.Status == "cerrado" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "proyecto cerrado"})
	}

	// permission check: only admin or project manager
	if v := c.Locals("user_id"); v == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no autorizado"})
	} else {
		var uid uint
		switch id := v.(type) {
		case uint:
			uid = id
		case int:
			uid = uint(id)
		case int64:
			uid = uint(id)
		case float64:
			uid = uint(id)
		}
		if !isProjectManager(pid, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "solo gerente o admin puede crear planes"})
		}
	}

	// validate responsable exists if provided
	if body.ResponsableID != nil {
		var u models.User
		if err := database.DB.First(&u, *body.ResponsableID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "responsable no encontrado"})
		}
	}

	if err := database.DB.Create(&body).Error; err != nil {
		fmt.Println("CreatePlanAccion: failed to create plan, db error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo crear plan"})
	}
	// log
	LogAction(c, "Planes de Acción", "create", fmt.Sprintf("plan_id=%d proyecto=%d", body.ID, body.ProyectoID))
	return c.Status(fiber.StatusCreated).JSON(body)
}

// ListPlanesPorProyecto lists plans for a project
func ListPlanesPorProyecto(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var pid uint
	fmt.Sscanf(projectID, "%d", &pid)
	var planes []models.PlanAccion
	database.DB.Where("proyecto_id = ?", pid).Preload("Responsable").Find(&planes)
	return c.JSON(planes)
}

// GetPlanAccion
func GetPlanAccion(c *fiber.Ctx) error {
	id := c.Params("id")
	var pl models.PlanAccion
	if err := database.DB.Preload("Responsable").First(&pl, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "plan no encontrado"})
	}
	return c.JSON(pl)
}

// UpdatePlanAccion
func UpdatePlanAccion(c *fiber.Ctx) error {
	id := c.Params("id")
	var req planReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	var pl models.PlanAccion
	if err := database.DB.First(&pl, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "plan no encontrado"})
	}
	// permission check
	var uid uint
	if v := c.Locals("user_id"); v == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no autorizado"})
	} else {
		switch idv := v.(type) {
		case uint:
			uid = idv
		case int:
			uid = uint(idv)
		case int64:
			uid = uint(idv)
		case float64:
			uid = uint(idv)
		}
		if !isProjectManager(pl.ProyectoID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "solo gerente o admin puede modificar planes"})
		}
	}
	// apply updates
	if req.Actividad != "" { pl.Actividad = req.Actividad }
	if req.Accion != "" { pl.Accion = req.Accion }
	if req.CantidadHoras != nil { pl.CantidadHoras = *req.CantidadHoras }
	// responsable can be set here, and when changed we must propagate to costos
	var responsableChanged bool
	if req.ResponsableID != nil {
		if pl.ResponsableID == nil || *req.ResponsableID != *pl.ResponsableID {
			pl.ResponsableID = req.ResponsableID
			responsableChanged = true
		}
	}
	// parse and set dates
	if req.FechaInicio != nil && *req.FechaInicio != "" {
		if t, err := time.Parse(time.RFC3339, *req.FechaInicio); err == nil {
			pl.FechaInicio = &t
		} else if t2, err2 := time.Parse("2006-01-02", *req.FechaInicio); err2 == nil {
			pl.FechaInicio = &t2
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fecha_inicio inválida"})
		}
	}
	if req.FechaCierre != nil && *req.FechaCierre != "" {
		if t, err := time.Parse(time.RFC3339, *req.FechaCierre); err == nil {
			pl.FechaCierre = &t
		} else if t2, err2 := time.Parse("2006-01-02", *req.FechaCierre); err2 == nil {
			pl.FechaCierre = &t2
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fecha_cierre inválida"})
		}
	}
	// validate dates
	if pl.FechaInicio != nil && pl.FechaCierre != nil {
		if pl.FechaCierre.Before(*pl.FechaInicio) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "fecha_cierre debe ser >= fecha_inicio"})
		}
	}
	if pl.CantidadHoras < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cantidad_horas inválida"})
	}
	database.DB.Save(&pl)
	// if responsable changed, update all costos to use the same responsable
	if responsableChanged {
		database.DB.Model(&models.CostoRecursoHumano{}).Where("plan_accion_id = ?", pl.ID).Update("responsable_id", pl.ResponsableID)
		database.DB.Model(&models.CostoMaterialInsumo{}).Where("plan_accion_id = ?", pl.ID).Update("responsable_id", pl.ResponsableID)
	}
	LogAction(c, "Planes de Acción", "update", fmt.Sprintf("plan_id=%d proyecto=%d", pl.ID, pl.ProyectoID))
	return c.JSON(pl)
}

// DeletePlanAccion
func DeletePlanAccion(c *fiber.Ctx) error {
	id := c.Params("id")
	var pl models.PlanAccion
	if err := database.DB.First(&pl, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "plan no encontrado"})
	}
	// permission check
	var uid uint
	if v := c.Locals("user_id"); v == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no autorizado"})
	} else {
		switch idv := v.(type) {
		case uint:
			uid = idv
		case int:
			uid = uint(idv)
		case int64:
			uid = uint(idv)
		case float64:
			uid = uint(idv)
		}
		if !isProjectManager(pl.ProyectoID, uid) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "solo gerente o admin puede eliminar planes"})
		}
	}
	if err := database.DB.Delete(&models.PlanAccion{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo eliminar plan"})
	}
	LogAction(c, "Planes de Acción", "delete", fmt.Sprintf("plan_id=%d proyecto=%d", pl.ID, pl.ProyectoID))
	return c.JSON(fiber.Map{"message": "plan eliminado"})
}
