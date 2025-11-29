package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/middleware"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(app *fiber.App) {
	// Auth routes v1 (public)
	authV1 := app.Group("/api/v1/auth")
	authV1.Post("/login", LoginHandler)
	authV1.Post("/register", RegisterHandler)
	authV1.Post("/refresh", RefreshTokenHandler)
	authV1.Post("/logout", LogoutHandler)

	// Protected routes v1 (require authentication)
	protectedV1 := app.Group("/api/v1/auth", middleware.AuthMiddleware)
	protectedV1.Get("/profile", GetProfileHandler)

	// Achievement routes v1 (require authentication)
	achievementV1 := app.Group("/api/v1/achievements", middleware.AuthMiddleware)

	// List (filtered by role) - read permission
	achievementV1.Get("", middleware.RBACMiddleware("achievement:read"), AchievementListHandler)

	// Detail - read permission
	achievementV1.Get("/:id", middleware.RBACMiddleware("achievement:read"), AchievementDetailHandler)

	// Create/Submit Prestasi (FR-003) - Mahasiswa & Admin
	achievementV1.Post("", middleware.RBACMiddleware("achievement:create"), AchievementCreateHandler)

	// Update - update permission
	achievementV1.Put("/:id", middleware.RBACMiddleware("achievement:update"), AchievementUpdateHandler)

	// Delete - delete permission
	achievementV1.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), AchievementDeleteHandler)

	// Submit for verification - submit permission
	achievementV1.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), AchievementSubmitHandler)

	// Verify (Dosen Wali & Admin) - verify permission
	achievementV1.Post("/:id/verify", middleware.RBACMiddleware("achievement:verify"), AchievementVerifyHandler)

	// Reject (Dosen Wali & Admin) - reject permission
	achievementV1.Post("/:id/reject", middleware.RBACMiddleware("achievement:reject"), AchievementRejectHandler)

	// Status history - read permission (semua role bisa lihat)
	achievementV1.Get("/:id/history", middleware.RBACMiddleware("achievement:read"), AchievementHistoryHandler)

	// Upload attachments - Mahasiswa & Admin
	achievementV1.Post("/:id/attachments", middleware.RBACMiddleware("achievement:upload"), AchievementUploadAttachmentHandler)
}
