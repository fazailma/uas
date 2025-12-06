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
	achievements.Get("/", middleware.RBACMiddleware("achievement:read"), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)

		achievements, err := achievementService.ListAchievements(userID, role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"code":   500,
				"error":  "failed to fetch achievements",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"code":   200,
			"data":   achievements,
		})
	})

	// GET /api/v1/achievements/:id - Detail
	achievements.Get("/:id", middleware.RBACMiddleware("achievement:read"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)

		achievement, err := achievementService.GetAchievementDetail(id, userID, role)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"code":   404,
				"error":  "achievement not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"code":   200,
			"data":   achievement,
		})
	})

	// POST /api/v1/achievements - Create (Mahasiswa)
	achievements.Post("/", middleware.RBACMiddleware("achievement:create"), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)

		var req struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
			Category    string `json:"category" binding:"required"`
			Date        string `json:"date" binding:"required"`
			ProofURL    string `json:"proof_url"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  "invalid request body",
			})
		}

		achievement, err := achievementService.CreateAchievement(userID, role, req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  "success",
			"code":    201,
			"message": "Prestasi berhasil dibuat",
			"data":    achievement,
		})
	})

	// PUT /api/v1/achievements/:id - Update (Mahasiswa)
	achievements.Put("/:id", middleware.RBACMiddleware("achievement:update"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)

		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Category    string `json:"category"`
			Date        string `json:"date"`
			ProofURL    string `json:"proof_url"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  "invalid request body",
			})
		}

		achievement, err := achievementService.UpdateAchievement(id, userID, role, req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"code":    200,
			"message": "Prestasi berhasil diperbarui",
			"data":    achievement,
		})
	})

	// DELETE /api/v1/achievements/:id - Delete (Mahasiswa)
	achievements.Delete("/:id", middleware.RBACMiddleware("achievement:delete"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)

		err := achievementService.DeleteAchievement(id, userID, role)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"code":    200,
			"message": "Prestasi berhasil dihapus",
		})
	})

	// POST /api/v1/achievements/:id/submit - Submit for verification
	achievements.Post("/:id/submit", middleware.RBACMiddleware("achievement:submit"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)

		err := achievementService.SubmitAchievement(id, userID, role)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"code":    200,
			"message": "Prestasi berhasil disubmit untuk verifikasi",
			"data": fiber.Map{
				"id":     id,
				"status": "submitted",
			},
		})
	})

	// POST /api/v1/achievements/:id/verify - Verify (Dosen Wali)
	achievements.Post("/:id/verify", middleware.RBACMiddleware("achievement:verify"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		dosenID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)

		if role != "Dosen Wali" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"code":   403,
				"error":  "only dosen wali can verify achievements",
			})
		}

		err := achievementService.VerifyAchievement(id, dosenID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"code":    200,
			"message": "Prestasi berhasil diverifikasi",
			"data": fiber.Map{
				"id":     id,
				"status": "verified",
			},
		})
	})

	// POST /api/v1/achievements/:id/reject - Reject (Dosen Wali)
	achievements.Post("/:id/reject", middleware.RBACMiddleware("achievement:verify"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		role := c.Locals("role").(string)

		if role != "Dosen Wali" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"code":   403,
				"error":  "only dosen wali can reject achievements",
			})
		}

		var req struct {
			RejectionNote string `json:"rejection_note" binding:"required"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  "invalid request body",
			})
		}

		err := achievementService.RejectAchievement(id, req.RejectionNote)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"code":    200,
			"message": "Prestasi berhasil ditolak",
			"data": fiber.Map{
				"id":     id,
				"status": "rejected",
			},
		})
	})

	// GET /api/v1/achievements/:id/history - Status history
	achievements.Get("/:id/history", middleware.RBACMiddleware("achievement:read"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)

		history, err := achievementService.GetAchievementHistory(id, userID, role)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": "error",
				"code":   404,
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"code":   200,
			"data":   history,
		})
	})

	// POST /api/v1/achievements/:id/attachments - Upload files
	achievements.Post("/:id/attachments", middleware.RBACMiddleware("achievement:update"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)

		if role != "Mahasiswa" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"code":   403,
				"error":  "only mahasiswa can upload attachments",
			})
		}

		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "error",
				"code":   400,
				"error":  "file is required",
			})
		}

		err = achievementService.ValidateAchievementOwnership(id, userID)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"code":   403,
				"error":  err.Error(),
			})
		}

		// TODO: Implement file upload logic
		_ = file

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"code":    200,
			"message": "File uploaded successfully",
			"data": fiber.Map{
				"achievement_id": id,
				"file_name":      file.Filename,
				"file_size":      file.Size,
			},
		})
	})
}
