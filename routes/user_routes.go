package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

// SetupUserRoutes sets up user management routes
func SetupUserRoutes(app *fiber.App) {
	svc := service.NewUserService()
	g := app.Group("/api/v1/users", middleware.AuthMiddleware, middleware.RBACMiddleware("admin:manage"))

	// User Management
	g.Get("/", svc.ListUsers)
	g.Get("/:id", svc.GetUserByID)
	g.Post("/", svc.CreateUser)
	g.Put("/:id", svc.UpdateUser)
	g.Delete("/:id", svc.DeleteUser)
	g.Put("/:id/role", svc.UpdateUserRole)
}
