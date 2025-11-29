package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/middleware"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(app *fiber.App) {
	// Auth routes v1 (public)
	authV1 := app.Group("/api/v1/auth")
	authV1.Post("/login", LoginHandler)
	authV1.Post("/register", RegisterHandler)
	authV1.Post("/refresh", RefreshTokenHandler)
	authV1.Post("/logout", LogoutHandler)

	// Protected routes v1 (require authentication)
	protectedV1 := app.Group("/api/v1/auth", middleware.AuthMiddleware)
	protectedV1.Get("/profile", GetProfileHandler)

	// Setup achievement routes
	SetupAchievementRoutes(app)
}
