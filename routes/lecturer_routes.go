package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

func SetupLecturerRoutes(app *fiber.App) {
	svc := service.NewLecturerService()
	g := app.Group("/api/v1/lecturers", middleware.AuthMiddleware)

	// Lecturer Management
	g.Get("/", middleware.RBACMiddleware("lecturer:read"), svc.ListLecturers)
	// Create/Update lecturer profile - Admin only (user:manage permission)
	g.Post("/", middleware.RBACMiddleware("user:manage"), svc.CreateLecturerProfile)
	g.Put("/:id", middleware.RBACMiddleware("user:manage"), svc.UpdateLecturerProfile)
	g.Get("/:id/advisees", middleware.RBACMiddleware("lecturer:read"), svc.GetAdvisees)
}
