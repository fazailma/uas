package routes

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all routes for the application
// @Summary Setup all API routes
// @Description Configure all authentication, achievement, user, student, lecturer and report routes
func SetupRoutes(app *fiber.App) {
	// Setup auth routes
	SetupAuthRoutes(app)

	// Setup achievement routes
	SetupAchievementRoutes(app)

	// Setup user management routes
	SetupUserRoutes(app)

	// Setup student routes
	SetupStudentRoutes(app)

	// Setup lecturer routes
	SetupLecturerRoutes(app)

	// Setup report and analytics routes
	SetupReportRoutes(app)
}
