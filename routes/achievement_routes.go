package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

func SetupAchievementRoutes(app *fiber.App) {
	achievementService := service.NewAchievementService()

	achievements := app.Group("/api/v1/achievements", middleware.AuthMiddleware)

	// List achievements
	achievements.Get("/", middleware.RBACMiddleware("achievement:read"), achievementService.AchievementListHandler)

	// Get achievement detail
	achievements.Get("/:id", middleware.RBACMiddleware("achievement:read"), achievementService.AchievementDetailHandler)

	// FR-003: Create/Submit achievement
	achievements.Post("/", middleware.RBACMiddleware("achievement:create"), achievementService.AchievementCreateHandler)

	// FR-004: Submit achievement for verification
	achievements.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), achievementService.AchievementSubmitHandler)

	// Update achievement
	achievements.Put("/:id", middleware.RBACMiddleware("achievement:update"), achievementService.AchievementUpdateHandler)

	// FR-005: Delete achievement
	achievements.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), achievementService.AchievementDeleteHandler)
}
