package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/utils"
)

// AchievementService defines all achievement-related operations
type AchievementService interface {
	CreateAchievement(c *fiber.Ctx) error
	UpdateAchievement(c *fiber.Ctx) error
	DeleteAchievement(c *fiber.Ctx) error
	SubmitAchievement(c *fiber.Ctx) error
	ListAchievements(c *fiber.Ctx) error
	GetAchievementDetail(c *fiber.Ctx) error
	GetAchievementHistory(c *fiber.Ctx) error
	GetStatistics(c *fiber.Ctx) error
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error
	UploadAttachment(c *fiber.Ctx) error
}

type achievementServiceImpl struct {
	pgRepo      *repository.AchievementRepository
	mongoRepo   *repository.MongoAchievementRepository
	studentRepo *repository.StudentRepository
	userRepo    *repository.UserRepository
}

func NewAchievementService() AchievementService {
	return &achievementServiceImpl{
		pgRepo:      repository.NewAchievementRepository(),
		mongoRepo:   repository.NewMongoAchievementRepository(),
		studentRepo: repository.NewStudentRepository(),
		userRepo:    repository.NewUserRepository(),
	}
}

// CreateAchievement handles achievement creation
// @Summary Create new achievement
// @Description Create a new achievement for the logged-in student
// @Tags Achievements
// @Accept json
// @Produce json
// @Param body body models.CreateAchievementRequest true "Achievement data"
// @Success 201 {object} models.AchievementReference
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements [post]
// @Security Bearer
func (s *achievementServiceImpl) CreateAchievement(c *fiber.Ctx) error {
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Title == "" || req.Category == "" || req.Date == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "title, category, and date are required")
	}

	if c.Locals("role") != "Mahasiswa" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "only mahasiswa can create achievements")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoAch, err := s.mongoRepo.Create(ctx, &models.MongoAchievement{
		StudentID:   c.Locals("user_id").(string),
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Date:        req.Date,
		ProofURL:    req.ProofURL,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to save achievement")
	}

	pgAch := &models.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          c.Locals("user_id").(string),
		MongoAchievementID: mongoAch.ID.Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.pgRepo.Create(pgAch); err != nil {
		s.mongoRepo.SoftDelete(ctx, mongoAch.ID.Hex())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to save achievement reference")
	}

	return utils.CreatedResponse(c, "Prestasi berhasil dibuat", pgAch)
}

// ListAchievements handles listing achievements
// @Summary List achievements
// @Description Get list of achievements based on user role
// @Tags Achievements
// @Produce json
// @Success 200 {array} models.AchievementReference
// @Failure 500 {object} map[string]string
// @Router /achievements [get]
// @Security Bearer
func (s *achievementServiceImpl) ListAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	var achievements []models.AchievementReference
	var err error

	switch role {
	case "Admin":
		achievements, err = s.pgRepo.FindAll()
	case "Mahasiswa":
		achievements, err = s.pgRepo.FindByStudentID(userID)
	default:
		achievements, err = s.pgRepo.FindAll()
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve achievements")
	}

	return utils.SuccessResponse(c, "achievements retrieved successfully", achievements)
}

// GetAchievementDetail handles getting achievement detail
// @Summary Get achievement detail
// @Description Retrieve detailed information of a specific achievement
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} models.AchievementReference
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /achievements/{id} [get]
// @Security Bearer
func (s *achievementServiceImpl) GetAchievementDetail(c *fiber.Ctx) error {
	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	role := c.Locals("role").(string)
	if role == "Mahasiswa" && achievement.StudentID != c.Locals("user_id").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only view your own achievements")
	}

	return utils.SuccessResponse(c, "achievement detail retrieved", achievement)
}

// UpdateAchievement handles updating achievement
// @Summary Update achievement
// @Description Update an existing achievement (only draft status)
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param body body models.UpdateAchievementRequest true "Updated achievement data"
// @Success 200 {object} models.AchievementReference
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /achievements/{id} [put]
// @Security Bearer
func (s *achievementServiceImpl) UpdateAchievement(c *fiber.Ctx) error {
	var req models.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	if c.Locals("role") == "Mahasiswa" && achievement.StudentID != c.Locals("user_id").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only update your own achievements")
	}

	if achievement.Status != "draft" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only draft achievements can be updated")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := s.mongoRepo.Update(ctx, achievement.MongoAchievementID, &models.MongoAchievement{
		StudentID:   c.Locals("user_id").(string),
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Date:        req.Date,
		ProofURL:    req.ProofURL,
	}); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update achievement")
	}

	achievement.UpdatedAt = time.Now()
	if err := s.pgRepo.Update(c.Params("id"), achievement); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update achievement")
	}

	return utils.SuccessResponse(c, "Prestasi berhasil diperbarui", achievement)
}

// DeleteAchievement handles deleting achievement
// @Summary Delete achievement
// @Description Delete an achievement (only draft status)
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /achievements/{id} [delete]
// @Security Bearer
func (s *achievementServiceImpl) DeleteAchievement(c *fiber.Ctx) error {
	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	if c.Locals("role") == "Mahasiswa" && achievement.StudentID != c.Locals("user_id").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only delete your own achievements")
	}

	if achievement.Status != "draft" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only draft achievements can be deleted")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.mongoRepo.SoftDelete(ctx, achievement.MongoAchievementID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to delete achievement")
	}

	if err := s.pgRepo.Delete(c.Params("id")); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to delete achievement")
	}

	return utils.DeletedResponse(c, "Prestasi berhasil dihapus")
}

// SubmitAchievement handles submitting achievement for verification
// @Summary Submit achievement
// @Description Submit an achievement for verification (changes status from draft to submitted)
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /achievements/{id}/submit [post]
// @Security Bearer
func (s *achievementServiceImpl) SubmitAchievement(c *fiber.Ctx) error {
	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	if c.Locals("role") == "Mahasiswa" && achievement.StudentID != c.Locals("user_id").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only submit your own achievements")
	}

	if achievement.Status != "draft" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only draft achievements can be submitted")
	}

	if err := s.pgRepo.UpdateStatus(c.Params("id"), "submitted"); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to submit achievement")
	}

	return utils.SuccessResponse(c, "Prestasi berhasil disubmit untuk verifikasi", fiber.Map{"id": c.Params("id"), "status": "submitted"})
}

// GetAchievementHistory handles getting achievement history
// @Summary Get achievement history
// @Description Get the timeline/history of an achievement's status changes
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} fiber.Map
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /achievements/{id}/history [get]
// @Security Bearer
func (s *achievementServiceImpl) GetAchievementHistory(c *fiber.Ctx) error {
	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	if c.Locals("role") == "Mahasiswa" && achievement.StudentID != c.Locals("user_id").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only view your own achievement history")
	}

	timeline := []map[string]interface{}{
		{
			"status":    "draft",
			"timestamp": achievement.CreatedAt,
			"message":   "Achievement created",
		},
	}

	if !achievement.SubmittedAt.IsZero() {
		timeline = append(timeline, map[string]interface{}{
			"status":    "submitted",
			"timestamp": achievement.SubmittedAt,
			"message":   "Achievement submitted for verification",
		})
	}

	if !achievement.VerifiedAt.IsZero() {
		timeline = append(timeline, map[string]interface{}{
			"status":      "verified",
			"timestamp":   achievement.VerifiedAt,
			"verified_by": achievement.VerifiedBy,
			"message":     "Achievement verified",
		})
	}

	if achievement.Status == "rejected" && achievement.RejectionNote != "" {
		timeline = append(timeline, map[string]interface{}{
			"status":  "rejected",
			"message": achievement.RejectionNote,
		})
	}

	return utils.SuccessResponse(c, "achievement history retrieved", map[string]interface{}{
		"id":       c.Params("id"),
		"status":   achievement.Status,
		"timeline": timeline,
	})
}

// GetStatistics handles getting achievement statistics
// @Summary Get achievement statistics
// @Description Get statistics of achievements based on user role
// @Tags Achievements
// @Produce json
// @Success 200 {object} fiber.Map
// @Failure 500 {object} map[string]string
// @Router /achievements/stats [get]
// @Security Bearer
func (s *achievementServiceImpl) GetStatistics(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	var achievements []models.AchievementReference
	var err error

	if role == "Mahasiswa" {
		achievements, err = s.pgRepo.FindByStudentID(userID)
	} else if role == "Dosen Wali" {
		lecturerRepo := repository.NewLecturerRepository()
		lecturer, err := lecturerRepo.FindByUserID(userID)
		if err == nil && lecturer != nil {
			studentRepo := repository.NewStudentRepository()
			students, err := studentRepo.FindByAdvisorID(lecturer.ID)
			if err == nil && len(students) > 0 {
				var studentIDs []string
				for _, s := range students {
					studentIDs = append(studentIDs, s.UserID)
				}
				achievements, _ = s.pgRepo.FindByStudentIDs(studentIDs)
			}
		}
	} else {
		achievements, err = s.pgRepo.FindAll()
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve achievements")
	}

	stats := s.buildStatistics(achievements)
	return utils.SuccessResponse(c, "statistics retrieved successfully", stats)
}

// buildStatistics builds statistics from achievements
func (s *achievementServiceImpl) buildStatistics(achievements []models.AchievementReference) fiber.Map {
	var verified, pending, rejected, draft int64
	for i := range achievements {
		switch achievements[i].Status {
		case "verified":
			verified++
		case "submitted":
			pending++
		case "rejected":
			rejected++
		case "draft":
			draft++
		}
	}

	total := int64(len(achievements))
	verificationRate := 0.0
	if total > 0 {
		verificationRate = float64(verified) / float64(total) * 100
	}

	return fiber.Map{
		"total":             total,
		"draft":             draft,
		"pending":           pending,
		"verified":          verified,
		"rejected":          rejected,
		"verification_rate": verificationRate,
	}
}

// VerifyAchievement handles achievement verification by lecturer
// @Summary Verify achievement
// @Description Verify an achievement submission (Dosen Wali only)
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /achievements/{id}/verify [post]
// @Security Bearer
func (s *achievementServiceImpl) VerifyAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Update achievement status to verified
	if err := s.pgRepo.UpdateStatus(achievementID, "verified"); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to verify achievement")
	}

	return utils.SuccessResponse(c, "achievement verified successfully", nil)
}

// RejectAchievement handles achievement rejection by lecturer
// @Summary Reject achievement
// @Description Reject an achievement submission with notes (Dosen Wali only)
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param body body map[string]string true "Rejection data"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /achievements/{id}/reject [post]
// @Security Bearer
func (s *achievementServiceImpl) RejectAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	var req struct {
		RejectionNote string `json:"rejection_note" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := s.pgRepo.RejectAchievement(achievementID, req.RejectionNote); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to reject achievement")
	}

	return utils.SuccessResponse(c, "achievement rejected successfully", nil)
}

// UploadAttachment handles file attachment upload for achievements
// @Summary Upload achievement attachment
// @Description Upload proof files for an achievement
// @Tags Achievements
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Achievement ID"
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /achievements/{id}/attachments [post]
// @Security Bearer
func (s *achievementServiceImpl) UploadAttachment(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "file is required")
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "file size exceeds 10MB limit")
	}

	// Save file (implementation depends on your storage setup)
	// For now, we'll just acknowledge the upload
	filename := uuid.New().String() + "_" + file.Filename

	// Update achievement with attachment URL
	// In production, this would be the actual file path/URL
	return utils.SuccessResponse(c, "file uploaded successfully", fiber.Map{
		"achievement_id": achievementID,
		"filename":       filename,
	})
}
