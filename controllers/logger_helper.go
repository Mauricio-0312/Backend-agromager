package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// LogAction creates a log entry. Called inline from controllers where events occur.
// The last parameter nonAuthUserID is optional and can be provided when the action
// happens before the request has an authenticated user (e.g. signup, login).
func LogAction(c *fiber.Ctx, module, event, details string, nonAuthUserID ...uint) {
	var userID *uint = nil
	if v := c.Locals("user_id"); v != nil {
		switch id := v.(type) {
		case uint:
			userID = &id
		case int:
			u := uint(id)
			userID = &u
		case int64:
			u := uint(id)
			userID = &u
		case float64:
			u := uint(id)
			userID = &u
		default:
			// ignore
		}
	}
	// if there's no authenticated user in the context, use provided optional id
	if userID == nil && len(nonAuthUserID) > 0 {
		u := nonAuthUserID[0]
		userID = &u
	}
	lg := models.Logger{
		UserID:  userID,
		Module:  module,
		Event:   event,
		Details: details,
	}
	if err := database.DB.Create(&lg).Error; err != nil {
		// cannot return error from here; just log to console
		fmt.Println("failed to create log:", err)
	}
}
