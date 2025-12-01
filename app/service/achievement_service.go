package service

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"UAS/app/models"
	"UAS/app/repository"
)

type AchievementService struct {
	pgRepo    *repository.AchievementRepository
	mongoRepo *repository.MongoAchievementRepository
}

func NewAchievementService() *AchievementService {
	return &AchievementService{
		pgRepo:    repository.NewAchievementRepository(),
		mongoRepo: repository.NewMongoAchievementRepository(),
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
		"message": "Prestasi berhasil dibuat",
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

// FR-007: Verify Prestasi
// VerifyAchievementHandler handles verifying achievement by dosen wali
// POST /api/v1/achievements/:id/verify
func (s *AchievementService) VerifyAchievementHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	dosenID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	// Only Dosen Wali can verify
	if role != "Dosen Wali" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"error":  "only dosen wali can verify achievements",
		})
	}

	err := s.VerifyAchievement(id, dosenID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Prestasi berhasil diverifikasi",
		"data": fiber.Map{
			"id":     id,
			"status": "verified",
		},
	})
}

// FR-008: Reject Prestasi
// RejectAchievementHandler handles rejecting achievement by dosen wali
// POST /api/v1/achievements/:id/reject
func (s *AchievementService) RejectAchievementHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role").(string)

	// Only Dosen Wali can reject
	if role != "Dosen Wali" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"error":  "only dosen wali can reject achievements",
		})
	}

	var req struct {
		RejectionNote string `json:"rejection_note" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "invalid request body",
		})
	}

	err := s.RejectAchievement(id, req.RejectionNote)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Prestasi berhasil ditolak",
		"data": fiber.Map{
			"id":     id,
			"status": "rejected",
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

// CreateAchievement creates achievement in both MongoDB and PostgreSQL
// FR-003: Submit Prestasi
func (s *AchievementService) CreateAchievement(userID, role string, req models.AchievementCreateRequest) (*models.AchievementReference, error) {
	// Validate required fields
	if req.Title == "" || req.Category == "" || req.Date == "" {
		return nil, errors.New("title, category, and date are required")
	}

	// Only Mahasiswa can create
	if role != "Mahasiswa" {
		return nil, errors.New("only mahasiswa can submit achievements")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Create MongoDB document with full achievement details
	mongoAchievement := &models.MongoAchievement{
		StudentID:   userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Date:        req.Date,
		ProofURL:    req.ProofURL,
	}

	mongoAch, err := s.mongoRepo.Create(ctx, mongoAchievement)
	if err != nil {
		return nil, errors.New("failed to save achievement to MongoDB: " + err.Error())
	}

	// 2. Create PostgreSQL reference
	pgAchievement := &models.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          userID,
		MongoAchievementID: mongoAch.ID.Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	err = s.pgRepo.Create(pgAchievement)
	if err != nil {
		// Rollback MongoDB if PostgreSQL fails
		s.mongoRepo.SoftDelete(ctx, mongoAch.ID.Hex())
		return nil, errors.New("failed to save achievement reference to PostgreSQL: " + err.Error())
	}

	return pgAchievement, nil
}

// SubmitAchievement updates achievement status to 'submitted'
// FR-004: Submit untuk Verifikasi
func (s *AchievementService) SubmitAchievement(id, userID, role string) error {
	achievement, err := s.pgRepo.FindByID(id)
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

	// Update status to submitted in PostgreSQL
	err = s.pgRepo.UpdateStatus(id, "submitted")
	if err != nil {
		return err
	}

	// TODO: Create notification for dosen wali

	return nil
}

// UpdateAchievement updates achievement details (only draft)
func (s *AchievementService) UpdateAchievement(id, userID, role string, req models.AchievementUpdateRequest) (*models.AchievementReference, error) {
	pgAchievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("achievement not found")
	}

	// Ownership check - only mahasiswa can update their own
	if role == "Mahasiswa" && pgAchievement.StudentID != userID {
		return nil, errors.New("you can only update your own achievements")
	}

	// Can only update draft achievements
	if pgAchievement.Status != "draft" {
		return nil, errors.New("only draft achievements can be updated")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update MongoDB document
	mongoID := pgAchievement.MongoAchievementID
	mongoAch := &models.MongoAchievement{
		StudentID:   userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Date:        req.Date,
		ProofURL:    req.ProofURL,
	}

	_, err = s.mongoRepo.Update(ctx, mongoID, mongoAch)
	if err != nil {
		return nil, errors.New("failed to update achievement in MongoDB: " + err.Error())
	}

	// Update PostgreSQL reference timestamp
	pgAchievement.UpdatedAt = time.Now()
	err = s.pgRepo.Update(id, pgAchievement)
	if err != nil {
		return nil, errors.New("failed to update achievement reference in PostgreSQL: " + err.Error())
	}

	return pgAchievement, nil
}

// DeleteAchievement soft deletes achievement from both MongoDB and marks PostgreSQL
// FR-005: Hapus Prestasi
func (s *AchievementService) DeleteAchievement(id, userID, role string) error {
	pgAchievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Ownership check - only mahasiswa can delete their own
	if role == "Mahasiswa" && pgAchievement.StudentID != userID {
		return errors.New("you can only delete your own achievements")
	}

	// Can only delete draft
	if pgAchievement.Status != "draft" {
		return errors.New("only draft achievements can be deleted")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Soft delete MongoDB document
	mongoID := pgAchievement.MongoAchievementID
	err = s.mongoRepo.SoftDelete(ctx, mongoID)
	if err != nil {
		return errors.New("failed to delete achievement from MongoDB: " + err.Error())
	}

	// 2. Soft delete PostgreSQL reference
	err = s.pgRepo.Delete(id)
	if err != nil {
		return errors.New("failed to delete achievement reference from PostgreSQL: " + err.Error())
	}

	return nil
}

// VerifyAchievement verifies an achievement by dosen wali
// FR-007: Verify Prestasi
func (s *AchievementService) VerifyAchievement(achievementID, dosenID string) error {
	// Get achievement
	achievement, err := s.pgRepo.FindByID(achievementID)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Verify ownership - check if dosen is advisor of this student
	student, err := repository.NewStudentRepository().FindByUserID(achievement.StudentID)
	if err != nil {
		return errors.New("student not found")
	}

	lecturer, err := repository.NewLecturerRepository().FindByUserID(dosenID)
	if err != nil {
		return errors.New("lecturer not found")
	}

	if student.AdvisorID != lecturer.ID {
		return errors.New("you can only verify achievements of your guided students")
	}

	// Check status
	if achievement.Status != "submitted" {
		return errors.New("only submitted achievements can be verified")
	}

	// Update status to verified
	err = s.pgRepo.VerifyAchievement(achievementID, dosenID)
	if err != nil {
		return errors.New("failed to verify achievement: " + err.Error())
	}

	// TODO: Create notification for student

	return nil
}

// RejectAchievement rejects an achievement with rejection note
// FR-008: Reject Prestasi
func (s *AchievementService) RejectAchievement(achievementID, rejectionNote string) error {
	// Get achievement
	achievement, err := s.pgRepo.FindByID(achievementID)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Check status
	if achievement.Status != "submitted" {
		return errors.New("only submitted achievements can be rejected")
	}

	// Update status to rejected with note
	err = s.pgRepo.RejectAchievement(achievementID, rejectionNote)
	if err != nil {
		return errors.New("failed to reject achievement: " + err.Error())
	}

	// TODO: Create notification for student with rejection reason

	return nil
}

// ListAchievements lists achievements based on user role
func (s *AchievementService) ListAchievements(userID, role string) ([]models.AchievementReference, error) {
	switch role {
	case "Admin":
		return s.pgRepo.FindAll()
	case "Mahasiswa":
		return s.pgRepo.FindByStudentID(userID)
	case "Dosen Wali":
		// TODO: Implement logic to find students under guidance
		return []models.AchievementReference{}, nil
	default:
		return s.pgRepo.FindAll()
	}
}

// GetAchievementDetail gets achievement details with ownership check
func (s *AchievementService) GetAchievementDetail(id, userID, role string) (*models.AchievementReference, error) {
	achievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Ownership check for Mahasiswa
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return nil, errors.New("you can only view your own achievements")
	}

	return achievement, nil
}

// AchievementHistoryHandler handles getting achievement status history
// GET /api/v1/achievements/:id/history
func (s *AchievementService) AchievementHistoryHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	history, err := s.GetAchievementHistory(id, userID, role)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   history,
	})
}

// AchievementUploadHandler handles uploading achievement attachments
// POST /api/v1/achievements/:id/attachments
func (s *AchievementService) AchievementUploadHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	// Only Mahasiswa can upload to their own achievements
	if role != "Mahasiswa" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"error":  "only mahasiswa can upload attachments",
		})
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  "file is required",
		})
	}

	// Validate ownership
	achievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error",
			"error":  "achievement not found",
		})
	}

	if achievement.StudentID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"error":  "you can only upload to your own achievements",
		})
	}

	// TODO: Implement file upload logic (S3, GCS, or local storage)
	// For now, just return success with file info
	_ = file // Use file variable to avoid unused error

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "File uploaded successfully",
		"data": fiber.Map{
			"achievement_id": id,
			"file_name":      file.Filename,
			"file_size":      file.Size,
		},
	})
}

// GetAchievementHistory retrieves achievement status history
func (s *AchievementService) GetAchievementHistory(id, userID, role string) (fiber.Map, error) {
	achievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("achievement not found")
	}

	// Ownership check for Mahasiswa
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return nil, errors.New("you can only view your own achievement history")
	}

	// Build history timeline
	history := fiber.Map{
		"id":     id,
		"status": achievement.Status,
		"timeline": []fiber.Map{
			{
				"status":    "draft",
				"timestamp": achievement.CreatedAt,
				"message":   "Achievement created",
			},
		},
	}

	// Add submitted event if applicable
	if !achievement.SubmittedAt.IsZero() {
		timeline := history["timeline"].([]fiber.Map)
		timeline = append(timeline, fiber.Map{
			"status":    "submitted",
			"timestamp": achievement.SubmittedAt,
			"message":   "Achievement submitted for verification",
		})
		history["timeline"] = timeline
	}

	// Add verified event if applicable
	if !achievement.VerifiedAt.IsZero() {
		timeline := history["timeline"].([]fiber.Map)
		timeline = append(timeline, fiber.Map{
			"status":      "verified",
			"timestamp":   achievement.VerifiedAt,
			"verified_by": achievement.VerifiedBy,
			"message":     "Achievement verified",
		})
		history["timeline"] = timeline
	}

	// Add rejected event if applicable
	if achievement.Status == "rejected" && achievement.RejectionNote != "" {
		timeline := history["timeline"].([]fiber.Map)
		timeline = append(timeline, fiber.Map{
			"status":  "rejected",
			"message": achievement.RejectionNote,
		})
		history["timeline"] = timeline
	}

	return history, nil
}
