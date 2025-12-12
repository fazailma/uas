package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

func SetupAuthRoutes(app *fiber.App) {
	svc := service.NewAuthService()
	g := app.Group("/api/v1/auth")

	g.Post("/login", svc.Login)
	g.Post("/logout", svc.Logout)
	g.Post("/refresh", svc.RefreshToken)

	protected := g.Group("", middleware.AuthMiddleware)
	protected.Get("/profile", svc.GetProfile)
}
