package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

// SetupReportRoutes sets up report and analytics routes
func SetupReportRoutes(app *fiber.App) {
	userSvc := service.NewUserService()
	g := app.Group("/api/v1/reports", middleware.AuthMiddleware)

	// Report and Analytics
	g.Get("/statistics", middleware.RBACMiddleware("report:read"), userSvc.GetAchievementStats)
	g.Get("/student/:id", middleware.RBACMiddleware("report:read"), userSvc.GetStudentAchievements)
}
