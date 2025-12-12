package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"UAS/app/repository"
	"UAS/utils"
)

// LecturerService defines all lecturer-related operations
type LecturerService interface {
	ListLecturers(c *fiber.Ctx) error
	GetAdvisees(c *fiber.Ctx) error
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error
	GetGuidedStudentsAchievements(c *fiber.Ctx) error
}

type lecturerServiceImpl struct {
	lecturerRepo    *repository.LecturerRepository
	userRepo        *repository.UserRepository
	achievementRepo *repository.AchievementRepository
	studentRepo     *repository.StudentRepository
	mongoRepo       *repository.MongoAchievementRepository
}

func NewLecturerService() LecturerService {
	return &lecturerServiceImpl{
		lecturerRepo:    repository.NewLecturerRepository(),
		userRepo:        repository.NewUserRepository(),
		achievementRepo: repository.NewAchievementRepository(),
		studentRepo:     repository.NewStudentRepository(),
		mongoRepo:       repository.NewMongoAchievementRepository(),
	}
}

// VerifyAchievement handles achievement verification
// @Summary Verify achievement
// @Description Verify a submitted achievement (Dosen Wali only)
// @Tags Lecturer
// @Produce json
// @Param id path string true \"Achievement ID\"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /lecturer/achievements/{id}/verify [post]
// @Security Bearer
func (s *lecturerServiceImpl) VerifyAchievement(c *fiber.Ctx) error {
	if c.Locals("role") != "Dosen Wali" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "only dosen wali can verify achievements")
	}

	achievementID := c.Params("id")
	dosenID := c.Locals("user_id").(string)

	achievement, err := s.achievementRepo.FindByID(achievementID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	student, err := s.studentRepo.FindByUserID(achievement.StudentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
	}

	lecturer, err := s.lecturerRepo.FindByUserID(dosenID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "lecturer not found")
	}

	if student.AdvisorID != lecturer.ID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only verify achievements of your guided students")
	}

	if achievement.Status != "submitted" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only submitted achievements can be verified")
	}

	if err := s.achievementRepo.VerifyAchievement(achievementID, dosenID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to verify achievement")
	}

	return utils.SuccessResponse(c, "Prestasi berhasil diverifikasi", fiber.Map{"id": achievementID, "status": "verified"})
}

// RejectAchievement handles achievement rejection
// @Summary Reject achievement
// @Description Reject a submitted achievement (Dosen Wali only)
// @Tags Lecturer
// @Accept json
// @Produce json
// @Param id path string true \"Achievement ID\"
// @Param body body map[string]string true \"Rejection note\"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /lecturer/achievements/{id}/reject [post]
// @Security Bearer
func (s *lecturerServiceImpl) RejectAchievement(c *fiber.Ctx) error {
	if c.Locals("role") != "Dosen Wali" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "only dosen wali can reject achievements")
	}

	var req struct {
		RejectionNote string `json:"rejection_note"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	achievementID := c.Params("id")

	achievement, err := s.achievementRepo.FindByID(achievementID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	if achievement.Status != "submitted" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only submitted achievements can be rejected")
	}

	if err := s.achievementRepo.RejectAchievement(achievementID, req.RejectionNote); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to reject achievement")
	}

	return utils.SuccessResponse(c, "Prestasi berhasil ditolak", fiber.Map{"id": achievementID, "status": "rejected"})
}

// GetGuidedStudentsAchievements handles getting guided students achievements
// @Summary Get guided students achievements
// @Description Get all achievements of students guided by current lecturer
// @Tags Lecturer
// @Produce json
// @Success 200 {array} models.AchievementReference
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lecturer/achievements [get]
// @Security Bearer
func (s *lecturerServiceImpl) GetGuidedStudentsAchievements(c *fiber.Ctx) error {
	if c.Locals("role") != "Dosen Wali" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "only dosen wali can access this endpoint")
	}

	dosenID := c.Locals("user_id").(string)

	lecturer, err := s.lecturerRepo.FindByUserID(dosenID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "lecturer not found")
	}

	students, err := s.studentRepo.FindByAdvisorID(lecturer.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch guided students")
	}

	if len(students) == 0 {
		return utils.SuccessResponse(c, "guided students achievements retrieved", []map[string]interface{}{})
	}

	var studentIDs []string
	for _, student := range students {
		studentIDs = append(studentIDs, student.UserID)
	}

	achievements, err := s.achievementRepo.FindByStudentIDs(studentIDs)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch achievements")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result []map[string]interface{}
	for i := range achievements {
		ach := &achievements[i]
		mongoAch, err := s.mongoRepo.FindByID(ctx, ach.MongoAchievementID)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":             ach.ID,
			"student_id":     ach.StudentID,
			"title":          mongoAch.Title,
			"description":    mongoAch.Description,
			"category":       mongoAch.Category,
			"date":           mongoAch.Date,
			"proof_url":      mongoAch.ProofURL,
			"status":         ach.Status,
			"submitted_at":   ach.SubmittedAt,
			"verified_at":    ach.VerifiedAt,
			"verified_by":    ach.VerifiedBy,
			"rejection_note": ach.RejectionNote,
			"created_at":     ach.CreatedAt,
		})
	}

	return utils.SuccessResponse(c, "guided students achievements retrieved", result)
}

// ListLecturers handles listing all lecturers
// @Summary List all lecturers
// @Description Get paginated list of all lecturers
// @Tags Lecturers
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /lecturers [get]
// @Security Bearer
func (s *lecturerServiceImpl) ListLecturers(c *fiber.Ctx) error {
	lecturers, err := s.lecturerRepo.FindAll()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to list lecturers")
	}

	return utils.SuccessResponse(c, "lecturers retrieved successfully", fiber.Map{
		"data": lecturers,
	})
}

// GetAdvisees handles getting lecturer's advisees
// @Summary Get lecturer's advisees
// @Description Get list of students guided by a lecturer
// @Tags Lecturers
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lecturers/{id}/advisees [get]
// @Security Bearer
func (s *lecturerServiceImpl) GetAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	// Verify lecturer exists
	_, err := s.lecturerRepo.FindByUserID(lecturerID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "lecturer not found")
	}

	// Get advisees (students with this lecturer as advisor)
	students, err := s.studentRepo.FindByAdvisorID(lecturerID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get advisees")
	}

	return utils.SuccessResponse(c, "advisees retrieved successfully", fiber.Map{
		"data": students,
	})
}
