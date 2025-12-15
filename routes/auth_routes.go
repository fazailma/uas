package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

// SetupAuthRoutes sets up authentication routes
// @Summary Setup authentication routes
// @Description Configure login, logout, refresh token and profile endpoints
func SetupAuthRoutes(app *fiber.App) {
	svc := service.NewAuthService()
	g := app.Group("/api/v1/auth")

	// @Summary User login
	// @Description Authenticate user with credentials
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /auth/login [post]
	g.Post("/login", svc.Login)

	// @Summary User logout
	// @Description Logout user
	// @Tags Auth
	// @Security Bearer
	// @Success 200 {object} map[string]interface{}
	// @Router /auth/logout [post]
	g.Post("/logout", svc.Logout)

	// @Summary Refresh token
	// @Description Refresh expired JWT token
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /auth/refresh [post]
	g.Post("/refresh", svc.RefreshToken)

	protected := g.Group("", middleware.AuthMiddleware)

	// @Summary Get user profile
	// @Description Get current authenticated user profile
	// @Tags Auth
	// @Security Bearer
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /auth/profile [get]
	protected.Get("/profile", svc.GetProfile)
}
