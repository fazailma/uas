package service

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"UAS/app/models"
	"UAS/app/repository"
)

type AchievementService struct {
	repo *repository.AchievementRepository
}

func NewAchievementService() *AchievementService {
	return &AchievementService{
		repo: repository.NewAchievementRepository(),
	}
}

// AchievementListHandler handles listing achievements
// GET /api/v1/achievements
func (s *AchievementService) AchievementListHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	achievements, err := s.ListAchievements(userID, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "failed to fetch achievements",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   achievements,
	})
}

// AchievementDetailHandler handles getting achievement detail
// GET /api/v1/achievements/:id
func (s *AchievementService) AchievementDetailHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	achievement, err := s.GetAchievementDetail(id, userID, role)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  "achievement not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   achievement,
	})
}

// AchievementCreateHandler handles creating/submitting achievement
// POST /api/v1/achievements
func (s *AchievementService) AchievementCreateHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	var req models.AchievementCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "invalid request body",
		})
	}

	achievement, err := s.CreateAchievement(userID, role, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Prestasi berhasil disubmit",
		"data":    achievement,
	})
}

// AchievementUpdateHandler handles updating achievement
// PUT /api/v1/achievements/:id
func (s *AchievementService) AchievementUpdateHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	var req models.AchievementUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "invalid request body",
		})
	}

	achievement, err := s.UpdateAchievement(id, userID, role, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement updated successfully",
		"data":    achievement,
	})
}

// AchievementDeleteHandler handles deleting achievement
// DELETE /api/v1/achievements/:id
func (s *AchievementService) AchievementDeleteHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	err := s.DeleteAchievement(id, userID, role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement deleted successfully",
	})
}

// AchievementSubmitHandler handles submitting achievement for verification
// POST /api/v1/achievements/:id/submit
func (s *AchievementService) AchievementSubmitHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	err := s.SubmitAchievement(id, userID, role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement submitted successfully",
		"data": fiber.Map{
			"id":     id,
			"status": "submitted",
		},
	})
}

// AchievementVerifyHandler handles verifying achievement
// POST /api/v1/achievements/:id/verify
func (s *AchievementService) AchievementVerifyHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	var req models.AchievementVerifyRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "invalid request body",
		})
	}

	err := s.VerifyAchievement(id, role, req.Points)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement verified successfully",
		"data": fiber.Map{
			"id":     id,
			"status": "verified",
			"points": req.Points,
		},
	})
}

// AchievementRejectHandler handles rejecting achievement
// POST /api/v1/achievements/:id/reject
func (s *AchievementService) AchievementRejectHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)
	var req models.AchievementRejectRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "invalid request body",
		})
	}

	err := s.RejectAchievement(id, role, req.Reason)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Achievement rejected",
		"data": fiber.Map{
			"id":     id,
			"status": "rejected",
			"reason": req.Reason,
		},
	})
}

// AchievementHistoryHandler handles getting achievement status history
// GET /api/v1/achievements/:id/history
func (s *AchievementService) AchievementHistoryHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	// TODO: Implement history tracking
	// For now, return empty history
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":      id,
			"history": []fiber.Map{},
		},
	})
}

// AchievementUploadAttachmentHandler handles uploading files to achievement
// POST /api/v1/achievements/:id/attachments
func (s *AchievementService) AchievementUploadAttachmentHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	// Only Mahasiswa and Admin can upload
	if role != "Mahasiswa" && role != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"error":  "only mahasiswa and admin can upload files",
		})
	}

	// Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "file is required",
		})
	}

	err = s.UploadAttachment(id, userID, role, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "File uploaded successfully",
		"data": fiber.Map{
			"id":       id,
			"filename": file.Filename,
			"size":     file.Size,
		},
	})
}

func (s *AchievementService) ListAchievements(userID, role string) ([]models.Achievement, error) {
	switch role {
	case "Admin":
		return s.repo.FindAll()
	case "Mahasiswa":
		return s.repo.FindByStudentID(userID)
	case "Dosen Wali":
		// TODO: Implement logic to find students under guidance
		return []models.Achievement{}, nil
	default:
		return s.repo.FindAll()
	}
}

func (s *AchievementService) GetAchievementDetail(id, userID, role string) (*models.Achievement, error) {
	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Ownership check for Mahasiswa
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return nil, errors.New("you can only view your own achievements")
	}

	// TODO: Dosen Wali check - verify mahasiswa is under their guidance
	if role == "Dosen Wali" {
		// For now allow all
	}

	return achievement, nil
}

func (s *AchievementService) CreateAchievement(userID, role string, req models.AchievementCreateRequest) (*models.Achievement, error) {
	// Validate required fields
	if req.Title == "" || req.Category == "" || req.Date == "" {
		return nil, errors.New("title, category, and date are required")
	}

	// Only Mahasiswa and Admin can create
	if role != "Mahasiswa" && role != "Admin" {
		return nil, errors.New("only mahasiswa and admin can submit achievements")
	}

	// Parse date
	achievementDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		achievementDate = time.Now()
	}

	achievement := &models.Achievement{
		ID:          uuid.New().String(),
		StudentID:   userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Date:        achievementDate,
		ProofURL:    req.ProofURL,
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.repo.Create(achievement)
	if err != nil {
		return nil, err
	}

	return achievement, nil
}

func (s *AchievementService) UpdateAchievement(id, userID, role string, req models.AchievementUpdateRequest) (*models.Achievement, error) {
	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("achievement not found")
	}

	// Ownership check
	if role == "Mahasiswa" {
		if achievement.StudentID != userID {
			return nil, errors.New("you can only update your own achievements")
		}
		// Mahasiswa hanya bisa update draft
		if achievement.Status != "draft" {
			return nil, errors.New("you can only update achievements in draft status")
		}
	}

	// Update fields
	if req.Title != "" {
		achievement.Title = req.Title
	}
	if req.Description != "" {
		achievement.Description = req.Description
	}
	if req.Category != "" {
		achievement.Category = req.Category
	}
	if req.Date != "" {
		if date, err := time.Parse("2006-01-02", req.Date); err == nil {
			achievement.Date = date
		}
	}
	if req.ProofURL != "" {
		achievement.ProofURL = req.ProofURL
	}

	achievement.UpdatedAt = time.Now()

	err = s.repo.Update(id, achievement)
	if err != nil {
		return nil, err
	}

	return achievement, nil
}

func (s *AchievementService) DeleteAchievement(id, userID, role string) error {
	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Ownership check
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return errors.New("you can only delete your own achievements")
	}

	return s.repo.Delete(id)
}

func (s *AchievementService) SubmitAchievement(id, userID, role string) error {
	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Ownership check
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return errors.New("you can only submit your own achievements")
	}

	// Can only submit draft
	if achievement.Status != "draft" {
		return errors.New("only draft achievements can be submitted")
	}

	return s.repo.UpdateStatus(id, "submitted")
}

func (s *AchievementService) VerifyAchievement(id, role string, points int) error {
	// Only Dosen Wali and Admin can verify
	if role != "Dosen Wali" && role != "Admin" {
		return errors.New("only dosen wali and admin can verify achievements")
	}

	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Can only verify submitted
	if achievement.Status != "submitted" {
		return errors.New("only submitted achievements can be verified")
	}

	achievement.Status = "verified"
	achievement.Points = points
	achievement.UpdatedAt = time.Now()

	return s.repo.Update(id, achievement)
}

func (s *AchievementService) RejectAchievement(id, role, reason string) error {
	// Only Dosen Wali and Admin can reject
	if role != "Dosen Wali" && role != "Admin" {
		return errors.New("only dosen wali and admin can reject achievements")
	}

	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Can only reject submitted
	if achievement.Status != "submitted" {
		return errors.New("only submitted achievements can be rejected")
	}

	achievement.Status = "rejected"
	achievement.RejectionNote = reason
	achievement.UpdatedAt = time.Now()

	return s.repo.Update(id, achievement)
}

func (s *AchievementService) UploadAttachment(id, userID, role, filename string) error {
	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Ownership check for Mahasiswa
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return errors.New("you can only upload files to your own achievements")
	}

	// TODO: Implement file upload to storage
	// For now, just update proof_url with filename
	achievement.ProofURL = filename
	achievement.UpdatedAt = time.Now()

	return s.repo.Update(id, achievement)
}
