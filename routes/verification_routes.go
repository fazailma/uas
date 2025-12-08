package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

func SetupVerificationRoutes(app *fiber.App) {
	svc := service.NewVerificationService()
	g := app.Group("/api/v1/verifications", middleware.AuthMiddleware)

	g.Get("/achievements", middleware.RBACMiddleware("achievement:read"), svc.GetAchievementsHandler)
	g.Post("/achievements/:id/verify", middleware.RBACMiddleware("achievement:verify"), svc.VerifyAchievementHandler)
	g.Post("/achievements/:id/reject", middleware.RBACMiddleware("achievement:verify"), svc.RejectAchievementHandler)
}
