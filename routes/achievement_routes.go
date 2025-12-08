package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

func SetupAchievementRoutes(app *fiber.App) {
	svc := service.NewAchievementService()
	g := app.Group("/api/v1/achievements", middleware.AuthMiddleware)

	g.Get("/", middleware.RBACMiddleware("achievement:read"), svc.ListAchievements)
	g.Get("/:id", middleware.RBACMiddleware("achievement:read"), svc.GetAchievementDetail)
	g.Post("/", middleware.RBACMiddleware("achievement:create"), svc.CreateAchievement)
	g.Put("/:id", middleware.RBACMiddleware("achievement:update"), svc.UpdateAchievement)
	g.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), svc.DeleteAchievement)
	g.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), svc.SubmitAchievement)
	g.Post("/:id/verify", middleware.RBACMiddleware("achievement:verify"), svc.VerifyAchievement)
	g.Post("/:id/reject", middleware.RBACMiddleware("achievement:verify"), svc.RejectAchievement)
	g.Get("/:id/history", middleware.RBACMiddleware("achievement:read"), svc.GetAchievementHistory)
	g.Post("/:id/attachments", middleware.RBACMiddleware("achievement:update"), svc.UploadAttachment)
}
