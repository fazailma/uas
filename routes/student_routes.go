package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

// SetupStudentRoutes sets up student management routes
func SetupStudentRoutes(app *fiber.App) {
	svc := service.NewStudentService()
	g := app.Group("/api/v1/students", middleware.AuthMiddleware)

	// Student Management
	g.Get("/", middleware.RBACMiddleware("student:read"), svc.ListStudents)
	g.Get("/:id", middleware.RBACMiddleware("student:read"), svc.GetStudent)
	g.Get("/:id/achievements", middleware.RBACMiddleware("achievement:read"), svc.GetStudentAchievements)
	g.Put("/:id/advisor", middleware.RBACMiddleware("student:update"), svc.SetAdvisor)
}
