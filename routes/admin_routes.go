package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

func SetupAdminRoutes(app *fiber.App) {
	svc := service.NewAdminService()
	g := app.Group("/api/v1/admin", middleware.AuthMiddleware, middleware.RBACMiddleware("admin:manage"))

	// User Management
	g.Get("/users", svc.ListUsers)
	g.Get("/users/:id", svc.GetUserByID)
	g.Post("/users", svc.CreateUser)
	g.Put("/users/:id", svc.UpdateUser)
	g.Delete("/users/:id", svc.DeleteUser)
	g.Post("/users/:student_id/set-advisor/:advisor_id", svc.SetAdvisor)

	// Achievement Management
	g.Get("/achievements", svc.GetAchievements)
	g.Get("/achievements/stats", svc.GetAchievementStats)
}
