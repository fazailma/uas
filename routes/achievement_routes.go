package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

func SetupAchievementRoutes(app *fiber.App) {
	achievementService := service.NewAchievementService()

	achievements := app.Group("/api/v1/achievements", middleware.AuthMiddleware)

	// GET /api/v1/achievements - List (filtered by role)
	achievements.Get("/", middleware.RBACMiddleware("achievement:read"), achievementService.AchievementListHandler)

	// GET /api/v1/achievements/:id - Detail
	achievements.Get("/:id", middleware.RBACMiddleware("achievement:read"), achievementService.AchievementDetailHandler)

	// POST /api/v1/achievements - Create (Mahasiswa)
	achievements.Post("/", middleware.RBACMiddleware("achievement:create"), achievementService.AchievementCreateHandler)

	// PUT /api/v1/achievements/:id - Update (Mahasiswa)
	achievements.Put("/:id", middleware.RBACMiddleware("achievement:update"), achievementService.AchievementUpdateHandler)

	// DELETE /api/v1/achievements/:id - Delete (Mahasiswa)
	achievements.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), achievementService.AchievementDeleteHandler)

	// POST /api/v1/achievements/:id/submit - Submit for verification
	achievements.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), achievementService.AchievementSubmitHandler)

	// POST /api/v1/achievements/:id/verify - Verify (Dosen Wali)
	achievements.Post("/:id/verify", middleware.RBACMiddleware("achievement:verify"), achievementService.VerifyAchievementHandler)

	// POST /api/v1/achievements/:id/reject - Reject (Dosen Wali)
	achievements.Post("/:id/reject", middleware.RBACMiddleware("achievement:verify"), achievementService.RejectAchievementHandler)

	// GET /api/v1/achievements/:id/history - Status history
	achievements.Get("/:id/history", middleware.RBACMiddleware("achievement:read"), achievementService.AchievementHistoryHandler)

	// POST /api/v1/achievements/:id/attachments - Upload files
	achievements.Post("/:id/attachments", middleware.RBACMiddleware("achievement:update"), achievementService.AchievementUploadHandler)
}
