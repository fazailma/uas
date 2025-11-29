package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

func SetupAchievementRoutes(app *fiber.App) {
	achievementService := service.NewAchievementService()

	achievements := app.Group("/api/v1/achievements", middleware.AuthMiddleware)

	achievements.Get("/", middleware.RBACMiddleware("achievement:read"), achievementService.AchievementListHandler)
	achievements.Get("/:id", middleware.RBACMiddleware("achievement:read"), achievementService.AchievementDetailHandler)
	achievements.Post("/", middleware.RBACMiddleware("achievement:create"), achievementService.AchievementCreateHandler)
	achievements.Put("/:id", middleware.RBACMiddleware("achievement:update"), achievementService.AchievementUpdateHandler)
	achievements.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), achievementService.AchievementDeleteHandler)
	achievements.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), achievementService.AchievementSubmitHandler)
	achievements.Post("/:id/verify", middleware.RBACMiddleware("achievement:verify"), achievementService.AchievementVerifyHandler)
	achievements.Post("/:id/reject", middleware.RBACMiddleware("achievement:reject"), achievementService.AchievementRejectHandler)
	achievements.Get("/:id/history", middleware.RBACMiddleware("achievement:read"), achievementService.AchievementHistoryHandler)
	achievements.Post("/:id/attachments", middleware.RBACMiddleware("achievement:upload"), achievementService.AchievementUploadAttachmentHandler)
}
