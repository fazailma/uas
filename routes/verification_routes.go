package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/middleware"
)

func SetupVerificationRoutes(app *fiber.App) {
	verificationService := service.NewVerificationService()

	verifications := app.Group("/api/v1/verifications", middleware.AuthMiddleware)

	// FR-006: Get achievements of guided students
	verifications.Get("/achievements", middleware.RBACMiddleware("achievement:read"), verificationService.ListGuidedStudentsAchievementsHandler)

	// FR-007: Verify achievement
	verifications.Post("/achievements/:id/verify", middleware.RBACMiddleware("achievement:verify"), verificationService.VerifyAchievementHandler)

	// FR-008: Reject achievement
	verifications.Post("/achievements/:id/reject", middleware.RBACMiddleware("achievement:verify"), verificationService.RejectAchievementHandler)
}
