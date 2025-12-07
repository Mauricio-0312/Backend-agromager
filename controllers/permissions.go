package controllers

import (
	"agroproject/backend/database"
	"agroproject/backend/models"
)

// isProjectManager returns true if user is admin or is a member of the project (user_projects table)
func isProjectManager(projectID uint, userID uint) bool {
	// admin shortcut
	var u models.User
	if err := database.DB.First(&u, userID).Error; err == nil {
		if u.Role == "admin" { return true }
	}
	// check membership in user_projects join table
	var cnt int64
	database.DB.Raw("SELECT COUNT(1) FROM user_projects WHERE project_id = ? AND user_id = ?", projectID, userID).Scan(&cnt)
	return cnt > 0
}
