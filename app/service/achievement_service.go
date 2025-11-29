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

// FR-003: Submit Prestasi
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

// FR-004: Submit untuk Verifikasi
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
		"message": "Prestasi berhasil disubmit untuk verifikasi",
		"data": fiber.Map{
			"id":     id,
			"status": "submitted",
		},
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
		"message": "Prestasi berhasil diperbarui",
		"data":    achievement,
	})
}

// FR-005: Hapus Prestasi
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
		"message": "Prestasi berhasil dihapus",
	})
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

// Business Logic Methods

func (s *AchievementService) CreateAchievement(userID, role string, req models.AchievementCreateRequest) (*models.AchievementReference, error) {
	// Validate required fields
	if req.Title == "" || req.Category == "" || req.Date == "" {
		return nil, errors.New("title, category, and date are required")
	}

	// Only Mahasiswa can create
	if role != "Mahasiswa" {
		return nil, errors.New("only mahasiswa can submit achievements")
	}

	achievement := &models.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          userID,
		MongoAchievementID: uuid.New().String(), // Will store MongoDB ID when data is saved to MongoDB
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	err := s.repo.Create(achievement)
	if err != nil {
		return nil, err
	}

	return achievement, nil
}

func (s *AchievementService) SubmitAchievement(id, userID, role string) error {
	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Ownership check - only mahasiswa can submit their own
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return errors.New("you can only submit your own achievements")
	}

	// Can only submit draft
	if achievement.Status != "draft" {
		return errors.New("only draft achievements can be submitted")
	}

	// Update status to submitted
	err = s.repo.UpdateStatus(id, "submitted")
	if err != nil {
		return err
	}

	// TODO: Create notification for dosen wali

	return nil
}

func (s *AchievementService) DeleteAchievement(id, userID, role string) error {
	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Ownership check - only mahasiswa can delete their own
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return errors.New("you can only delete your own achievements")
	}

	// Can only delete draft
	if achievement.Status != "draft" {
		return errors.New("only draft achievements can be deleted")
	}

	// Soft delete
	return s.repo.Delete(id)
}

func (s *AchievementService) ListAchievements(userID, role string) ([]models.AchievementReference, error) {
	switch role {
	case "Admin":
		return s.repo.FindAll()
	case "Mahasiswa":
		return s.repo.FindByStudentID(userID)
	case "Dosen Wali":
		// TODO: Implement logic to find students under guidance
		return []models.AchievementReference{}, nil
	default:
		return s.repo.FindAll()
	}
}

func (s *AchievementService) GetAchievementDetail(id, userID, role string) (*models.AchievementReference, error) {
	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Ownership check for Mahasiswa
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return nil, errors.New("you can only view your own achievements")
	}

	return achievement, nil
}

func (s *AchievementService) UpdateAchievement(id, userID, role string, req models.AchievementUpdateRequest) (*models.AchievementReference, error) {
	achievement, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("achievement not found")
	}

	// Ownership check - only mahasiswa can update their own
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return nil, errors.New("you can only update your own achievements")
	}

	// Can only update draft achievements
	if achievement.Status != "draft" {
		return nil, errors.New("only draft achievements can be updated")
	}

	// Update fields if provided
	if req.Title != "" {
		// Use the existing achievement data, we'll store it in MongoDB later
		// For now, we just mark it as updated
	}

	achievement.UpdatedAt = time.Now()

	err = s.repo.Update(id, achievement)
	if err != nil {
		return nil, err
	}

	return achievement, nil
}
