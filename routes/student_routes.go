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
	// Create/Update student profile - Admin only (user:manage permission)
	g.Post("/", middleware.RBACMiddleware("user:manage"), svc.CreateStudentProfile)
	g.Get("/:id", middleware.RBACMiddleware("student:read"), svc.GetStudent)
	g.Put("/:id", middleware.RBACMiddleware("user:manage"), svc.UpdateStudentProfile)
	g.Get("/:id/achievements", middleware.RBACMiddleware("achievement:read"), svc.GetStudentAchievements)
	g.Put("/:id/advisor", middleware.RBACMiddleware("user:manage"), svc.SetAdvisor)
}
