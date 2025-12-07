package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// CreateCostoMaterial
func CreateCostoMaterial(c *fiber.Ctx) error {
	planID := c.Params("id")
	var body models.CostoMaterialInsumo
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	var pid uint
	if _, err := fmt.Sscanf(planID, "%d", &pid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "plan id inválido"})
	}
	// validate plan exists
	var pl models.PlanAccion
	if err := database.DB.First(&pl, pid).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "plan no encontrado"})
	}
	// validations
	if body.Cantidad <= 0 || body.Costo <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cantidad y costo deben ser > 0"})
	}
	// unidad exists if provided
	if body.UnidadID != nil {
		var u models.UnitOfMeasure
		if err := database.DB.First(&u, *body.UnidadID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unidad no encontrada"})
		}
	}
	// use plan's responsable (do not allow override)
	body.ResponsableID = pl.ResponsableID
	body.PlanAccionID = pid
	// ensure action matches plan and calculate monto (read-only)
	body.Accion = pl.Accion
	body.Monto = body.Cantidad * body.Costo
	if err := database.DB.Create(&body).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo crear costo material"})
	}
	// update plan monto
	recalcPlanMonto(pid)
	LogAction(c, "Planes de Acción", "agregar costo material", fmt.Sprintf("costo_id=%d plan_id=%d monto=%.2f", body.ID, body.PlanAccionID, body.Monto))
	return c.Status(fiber.StatusCreated).JSON(body)
}

// List costos materiales por plan
func ListCostosMateriales(c *fiber.Ctx) error {
	planID := c.Params("id")
	var pid uint
	fmt.Sscanf(planID, "%d", &pid)
	var costos []models.CostoMaterialInsumo
	database.DB.Where("plan_accion_id = ?", pid).Preload("Unidad").Preload("Responsable").Find(&costos)
	return c.JSON(costos)
}

// GetCostoMaterial
func GetCostoMaterial(c *fiber.Ctx) error {
	id := c.Params("id")
	var cm models.CostoMaterialInsumo
	if err := database.DB.Preload("Unidad").Preload("Responsable").First(&cm, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "costo no encontrado"})
	}
	return c.JSON(cm)
}

// UpdateCostoMaterial
func UpdateCostoMaterial(c *fiber.Ctx) error {
	id := c.Params("id")
	var body models.CostoMaterialInsumo
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	var cm models.CostoMaterialInsumo
	if err := database.DB.First(&cm, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "costo no encontrado"})
	}
	if body.Cantidad != 0 {
		cm.Cantidad = body.Cantidad
	}
	if body.Costo != 0 {
		cm.Costo = body.Costo
	}
	if body.Actividad != "" {
		cm.Actividad = body.Actividad
	}
	// do NOT allow changing Accion or Responsable here; they are set from the PlanAccion
	if body.Categoria != 0 {
		cm.Categoria = body.Categoria
	}
	if body.Descripcion != "" {
		cm.Descripcion = body.Descripcion
	}
	if body.UnidadID != nil {
		cm.UnidadID = body.UnidadID
	}
	// responsable cannot be changed
	cm.Monto = cm.Cantidad * cm.Costo
	database.DB.Save(&cm)
	recalcPlanMonto(cm.PlanAccionID)
	LogAction(c, "Planes de Acción", "modificar costo material", fmt.Sprintf("costo_id=%d plan_id=%d monto=%.2f", cm.ID, cm.PlanAccionID, cm.Monto))
	return c.JSON(cm)
}

// DeleteCostoMaterial
func DeleteCostoMaterial(c *fiber.Ctx) error {
	id := c.Params("id")
	var cm models.CostoMaterialInsumo
	if err := database.DB.First(&cm, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "costo no encontrado"})
	}
	if err := database.DB.Delete(&models.CostoMaterialInsumo{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo eliminar costo"})
	}
	recalcPlanMonto(cm.PlanAccionID)
	LogAction(c, "Planes de Acción", "eliminar costo material", fmt.Sprintf("costo_id=%d plan_id=%d", cm.ID, cm.PlanAccionID))
	return c.JSON(fiber.Map{"message": "costo eliminado"})
}
