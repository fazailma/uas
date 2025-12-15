package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

// SetupReportRoutes sets up report and analytics routes
func SetupReportRoutes(app *fiber.App) {
	achievementSvc := service.NewAchievementService()
	g := app.Group("/api/v1/reports", middleware.AuthMiddleware)

	// Report and Analytics
	// Authorization handled in service layer: Admin sees all, Mahasiswa sees own, Dosen/Dosen Wali sees advisees
	g.Get("/statistics", achievementSvc.GetStatistics)
	g.Get("/student/:id", achievementSvc.GetStudentReport)
}
