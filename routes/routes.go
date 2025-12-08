package routes

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(app *fiber.App) {
	// Setup auth routes
	SetupAuthRoutes(app)

	// Setup achievement routes
	SetupAchievementRoutes(app)

	// Setup verification routes (for dosen wali)
	SetupVerificationRoutes(app)

	// Setup admin routes
	SetupAdminRoutes(app)
}
