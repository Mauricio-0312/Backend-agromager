package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateLog adds a log entry
func CreateLog(c *fiber.Ctx) error {
	var logm models.Logger
	if err := c.BodyParser(&logm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "input inválido"})
	}
	if err := database.DB.Create(&logm).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo crear log"})
	}
	return c.JSON(logm)
}

// List logs with optional query and time filters
func ListLogs(c *fiber.Ctx) error {
	var logs []models.Logger
	q := c.Query("q")
	// build base query
	db := database.DB.Model(&models.Logger{})
	filterApplied := "none"

	// Apply text search if provided
	if q != "" {
		db = db.Where("module LIKE ? OR event LIKE ? OR details LIKE ?", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	// Time filters (priority order): day, date, month, quarter, year, range
	// day/date -> YYYY-MM-DD (use DATE(created_at) = ?)
	// month+year -> month int, year int
	// quarter+year -> quarter 1-4, year int
	// year -> year int
	// range -> start & end YYYY-MM-DD

	day := c.Query("day")
	date := c.Query("date")
	month := c.Query("month")
	year := c.Query("year")
	quarter := c.Query("quarter")
	start := c.Query("start")
	end := c.Query("end")

	// determine which filter to apply based on priority
	// We'll compute start/end time ranges for each filter to use BETWEEN queries (portable)
	if day != "" {
		d, err := time.Parse("2006-01-02", day)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "day inválido, use YYYY-MM-DD"})
		}
		startT := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		endT := startT.Add(24*time.Hour - time.Nanosecond)
		db = db.Where("created_at BETWEEN ? AND ?", startT, endT)
		filterApplied = "day=" + day
	} else if date != "" {
		d, err := time.Parse("2006-01-02", date)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "date inválida, use YYYY-MM-DD"})
		}
		startT := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		endT := startT.Add(24*time.Hour - time.Nanosecond)
		db = db.Where("created_at BETWEEN ? AND ?", startT, endT)
		filterApplied = "date=" + date
	} else if month != "" && year != "" {
		m, errm := strconv.Atoi(month)
		y, erry := strconv.Atoi(year)
		if errm != nil || erry != nil || m < 1 || m > 12 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "month o year inválidos"})
		}
		startT := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		// compute first day of next month then subtract a nanosecond
		next := startT.AddDate(0, 1, 0)
		endT := next.Add(-time.Nanosecond)
		db = db.Where("created_at BETWEEN ? AND ?", startT, endT)
		filterApplied = fmt.Sprintf("month=%02d year=%04d", m, y)
	} else if quarter != "" && year != "" {
		qv, errq := strconv.Atoi(quarter)
		y, erry := strconv.Atoi(year)
		if errq != nil || erry != nil || qv < 1 || qv > 4 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "quarter o year inválidos"})
		}
		startMonth := (qv-1)*3 + 1
		startT := time.Date(y, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
		endT := startT.AddDate(0, 3, 0).Add(-time.Nanosecond)
		db = db.Where("created_at BETWEEN ? AND ?", startT, endT)
		filterApplied = fmt.Sprintf("quarter=%d year=%04d", qv, y)
	} else if year != "" {
		y, erry := strconv.Atoi(year)
		if erry != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "year inválido"})
		}
		startT := time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)
		endT := startT.AddDate(1, 0, 0).Add(-time.Nanosecond)
		db = db.Where("created_at BETWEEN ? AND ?", startT, endT)
		filterApplied = fmt.Sprintf("year=%04d", y)
	} else if start != "" && end != "" {
		startT, err1 := time.Parse("2006-01-02", start)
		endT, err2 := time.Parse("2006-01-02", end)
		if err1 != nil || err2 != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "start o end inválidos, usar YYYY-MM-DD"})
		}
		endT = endT.Add(24*time.Hour - time.Nanosecond)
		db = db.Where("created_at BETWEEN ? AND ?", startT, endT)
		filterApplied = "range=" + start + " to " + end
	}

	// Role based filter: non-admins see only their logs
	if v := c.Locals("user_id"); v != nil {
		var uid uint
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
		// fetch user role
		var u models.User
		if err := database.DB.First(&u, uid).Error; err == nil {
			if u.Role != "admin" {
				db = db.Where("user_id = ?", uid)
			}
		}
	}

	// preload user and order by creation time descending
	// debug: show applied filter
	fmt.Println("ListLogs: filterApplied=", filterApplied)
	res := db.Preload("User").Order("created_at desc").Find(&logs)
	if res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": res.Error.Error()})
	}
	// set header to help frontend debug what filter was applied
	c.Set("X-Filter-Applied", filterApplied)
	return c.JSON(logs)
}

// Get single log
func GetLog(c *fiber.Ctx) error {
	id := c.Params("id")
	var l models.Logger
	if err := database.DB.Preload("User").First(&l, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "log no encontrado"})
	}
	return c.JSON(l)
}

// Delete log
func DeleteLog(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.Logger{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo eliminar log"})
	}
	return c.JSON(fiber.Map{"message": "log eliminado"})
}

// Count logs
func CountLogs(c *fiber.Ctx) error {
	var count int64
	database.DB.Model(&models.Logger{}).Count(&count)
	return c.JSON(fiber.Map{"count": strconv.FormatInt(count, 10)})
}
