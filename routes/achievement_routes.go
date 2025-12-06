package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/models"
	"UAS/app/service"
	"UAS/helpers"
	"UAS/middleware"
)

func SetupAchievementRoutes(app *fiber.App) {
	svc := service.NewAchievementService()
	g := app.Group("/api/v1/achievements", middleware.AuthMiddleware)

	// GET / - List achievements
	g.Get("/", middleware.RBACMiddleware("achievement:read"), func(c *fiber.Ctx) error {
		data, err := svc.ListAchievements(c.Locals("user_id").(string), c.Locals("role").(string))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.BuildErrorResponse(500, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, data))
	})

	// GET /:id - Get achievement detail
	g.Get("/:id", middleware.RBACMiddleware("achievement:read"), func(c *fiber.Ctx) error {
		data, err := svc.GetAchievementDetail(c.Params("id"), c.Locals("user_id").(string), c.Locals("role").(string))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(helpers.BuildErrorResponse(404, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, data))
	})

	// POST / - Create achievement
	g.Post("/", middleware.RBACMiddleware("achievement:create"), func(c *fiber.Ctx) error {
		var req models.AchievementCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "invalid request body"))
		}
		data, err := svc.CreateAchievement(c.Locals("user_id").(string), c.Locals("role").(string), req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.BuildCreatedResponse("Prestasi berhasil dibuat", data))
	})

	// PUT /:id - Update achievement
	g.Put("/:id", middleware.RBACMiddleware("achievement:update"), func(c *fiber.Ctx) error {
		var req models.AchievementUpdateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "invalid request body"))
		}
		data, err := svc.UpdateAchievement(c.Params("id"), c.Locals("user_id").(string), c.Locals("role").(string), req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("Prestasi berhasil diperbarui", data))
	})

	// DELETE /:id - Delete achievement
	g.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), func(c *fiber.Ctx) error {
		if err := svc.DeleteAchievement(c.Params("id"), c.Locals("user_id").(string), c.Locals("role").(string)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildDeletedResponse("Prestasi berhasil dihapus"))
	})

	// POST /:id/submit - Submit achievement
	g.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), func(c *fiber.Ctx) error {
		if err := svc.SubmitAchievement(c.Params("id"), c.Locals("user_id").(string), c.Locals("role").(string)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("Prestasi berhasil disubmit untuk verifikasi", fiber.Map{"id": c.Params("id"), "status": "submitted"}))
	})

	// POST /:id/verify - Verify achievement
	g.Post("/:id/verify", middleware.RBACMiddleware("achievement:verify"), func(c *fiber.Ctx) error {
		if err := svc.VerifyAchievement(c.Params("id"), c.Locals("user_id").(string)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("Prestasi berhasil diverifikasi", fiber.Map{"id": c.Params("id"), "status": "verified"}))
	})

	// POST /:id/reject - Reject achievement
	g.Post("/:id/reject", middleware.RBACMiddleware("achievement:verify"), func(c *fiber.Ctx) error {
		var req struct {
			RejectionNote string `json:"rejection_note"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "invalid request body"))
		}
		if err := svc.RejectAchievement(c.Params("id"), req.RejectionNote); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("Prestasi berhasil ditolak", fiber.Map{"id": c.Params("id"), "status": "rejected"}))
	})

	// GET /:id/history - Get achievement history
	g.Get("/:id/history", middleware.RBACMiddleware("achievement:read"), func(c *fiber.Ctx) error {
		data, err := svc.GetAchievementHistory(c.Params("id"), c.Locals("user_id").(string), c.Locals("role").(string))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(helpers.BuildErrorResponse(404, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, data))
	})

	// POST /:id/attachments - Upload attachment
	g.Post("/:id/attachments", middleware.RBACMiddleware("achievement:update"), func(c *fiber.Ctx) error {
		if err := svc.ValidateAchievementOwnership(c.Params("id"), c.Locals("user_id").(string)); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(helpers.BuildErrorResponse(403, err.Error()))
		}
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "file is required"))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("File uploaded successfully", fiber.Map{"achievement_id": c.Params("id"), "file_name": file.Filename, "file_size": file.Size}))
	})
}
