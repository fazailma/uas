package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/service"
	"UAS/helpers"
	"UAS/middleware"
)

func SetupVerificationRoutes(app *fiber.App) {
	svc := service.NewVerificationService()
	g := app.Group("/api/v1/verifications", middleware.AuthMiddleware)

	// GET /achievements - Get guided students achievements
	g.Get("/achievements", middleware.RBACMiddleware("achievement:read"), func(c *fiber.Ctx) error {
		achievements, err := svc.GetGuidedStudentsAchievementsWithRoleCheck(c.Locals("user_id").(string), c.Locals("role").(string))
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(helpers.BuildErrorResponse(403, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, achievements))
	})

	// POST /achievements/:id/verify - Verify achievement
	g.Post("/achievements/:id/verify", middleware.RBACMiddleware("achievement:verify"), func(c *fiber.Ctx) error {
		if err := svc.VerifyAchievementWithRoleCheck(c.Params("id"), c.Locals("user_id").(string), c.Locals("role").(string)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("Prestasi berhasil diverifikasi", fiber.Map{"id": c.Params("id"), "status": "verified"}))
	})

	// POST /achievements/:id/reject - Reject achievement
	g.Post("/achievements/:id/reject", middleware.RBACMiddleware("achievement:verify"), func(c *fiber.Ctx) error {
		var req struct {
			RejectionNote string `json:"rejection_note"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "invalid request body"))
		}
		if err := svc.RejectAchievementWithRoleCheck(c.Params("id"), req.RejectionNote, c.Locals("role").(string)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("Prestasi berhasil ditolak", fiber.Map{"id": c.Params("id"), "status": "rejected"}))
	})
}
