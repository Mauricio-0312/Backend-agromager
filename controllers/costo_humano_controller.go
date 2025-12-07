package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// CreateCostoHumano
func CreateCostoHumano(c *fiber.Ctx) error {
	planID := c.Params("id")
	var body models.CostoRecursoHumano
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
	if body.Tiempo <= 0 || body.Cantidad <= 0 || body.Costo <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "tiempo, cantidad y costo deben ser > 0"})
	}
	// responsable exists: we'll ignore provided responsable and use plan's responsable
	//    if body.ResponsableID != nil {
	//        var u models.User
	//        if err := database.DB.First(&u, *body.ResponsableID).Error; err != nil {
	//            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "responsable no encontrado"})
	//        }
	//    }
	// use plan's responsable (do not allow override)
	body.ResponsableID = pl.ResponsableID
	body.PlanAccionID = pid
	// ensure action matches plan
	body.Accion = pl.Accion
	// calculate monto (read-only)
	body.Monto = body.Tiempo * body.Costo * float64(body.Cantidad)
	if err := database.DB.Create(&body).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo crear costo humano"})
	}
	// update plan monto (sum both types)
	recalcPlanMonto(pid)
	LogAction(c, "Planes de Acción", "agregar costo humano", fmt.Sprintf("costo_id=%d plan_id=%d monto=%.2f", body.ID, body.PlanAccionID, body.Monto))
	return c.Status(fiber.StatusCreated).JSON(body)
}

// List costos humanos por plan
func ListCostosHumanos(c *fiber.Ctx) error {
	planID := c.Params("id")
	var pid uint
	fmt.Sscanf(planID, "%d", &pid)
	var costos []models.CostoRecursoHumano
	database.DB.Where("plan_accion_id = ?", pid).Preload("Responsable").Find(&costos)
	return c.JSON(costos)
}

// GetCostoHumano
func GetCostoHumano(c *fiber.Ctx) error {
	id := c.Params("id")
	var ch models.CostoRecursoHumano
	if err := database.DB.Preload("Responsable").First(&ch, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "costo no encontrado"})
	}
	return c.JSON(ch)
}

// UpdateCostoHumano
func UpdateCostoHumano(c *fiber.Ctx) error {
	id := c.Params("id")
	var body models.CostoRecursoHumano
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	var ch models.CostoRecursoHumano
	if err := database.DB.First(&ch, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "costo no encontrado"})
	}
	// simple updates and validations
	if body.Tiempo != 0 {
		ch.Tiempo = body.Tiempo
	}
	if body.Cantidad != 0 {
		ch.Cantidad = body.Cantidad
	}
	if body.Costo != 0 {
		ch.Costo = body.Costo
	}
	if body.Actividad != "" {
		ch.Actividad = body.Actividad
	}
	// do NOT allow changing Accion or Responsable here; they are set from the PlanAccion
	// recalc monto
	ch.Monto = ch.Tiempo * ch.Costo * float64(ch.Cantidad)
	database.DB.Save(&ch)
	recalcPlanMonto(ch.PlanAccionID)
	LogAction(c, "Planes de Acción", "modificar costo humano", fmt.Sprintf("costo_id=%d plan_id=%d monto=%.2f", ch.ID, ch.PlanAccionID, ch.Monto))
	return c.JSON(ch)
}

// DeleteCostoHumano
func DeleteCostoHumano(c *fiber.Ctx) error {
	id := c.Params("id")
	var ch models.CostoRecursoHumano
	if err := database.DB.First(&ch, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "costo no encontrado"})
	}
	if err := database.DB.Delete(&models.CostoRecursoHumano{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo eliminar costo"})
	}
	recalcPlanMonto(ch.PlanAccionID)
	LogAction(c, "Planes de Acción", "eliminar costo humano", fmt.Sprintf("costo_id=%d plan_id=%d", ch.ID, ch.PlanAccionID))
	return c.JSON(fiber.Map{"message": "costo eliminado"})
}

// helper to recalc plan monto
func recalcPlanMonto(planID uint) {
	var sumH float64
	database.DB.Model(&models.CostoRecursoHumano{}).Where("plan_accion_id = ?", planID).Select("COALESCE(SUM(monto),0)").Row().Scan(&sumH)
	var sumM float64
	database.DB.Model(&models.CostoMaterialInsumo{}).Where("plan_accion_id = ?", planID).Select("COALESCE(SUM(monto),0)").Row().Scan(&sumM)
	var pl models.PlanAccion
	if err := database.DB.First(&pl, planID).Error; err != nil {
		return
	}
	pl.Monto = sumH + sumM
	database.DB.Save(&pl)
}
