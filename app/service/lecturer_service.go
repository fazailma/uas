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

// LecturerService defines all lecturer-related operations
type LecturerService interface {
	CreateLecturerProfile(c *fiber.Ctx) error
	ListLecturers(c *fiber.Ctx) error
	UpdateLecturerProfile(c *fiber.Ctx) error
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
	roleRepo        *repository.RoleRepository
}

func NewLecturerService() LecturerService {
	return &lecturerServiceImpl{
		lecturerRepo:    repository.NewLecturerRepository(),
		userRepo:        repository.NewUserRepository(),
		achievementRepo: repository.NewAchievementRepository(),
		studentRepo:     repository.NewStudentRepository(),
		mongoRepo:       repository.NewMongoAchievementRepository(),
		roleRepo:        repository.NewRoleRepository(),
	}
}

// CreateLecturerProfile handles creating lecturer profile
// @Summary Create lecturer profile
// @Description Create a new lecturer profile for a user
// @Tags Lecturers
// @Accept json
// @Produce json
// @Param body body models.Lecturer true "Lecturer data"
// @Success 201 {object} models.Lecturer
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lecturers [post]
// @Security Bearer
func (s *lecturerServiceImpl) CreateLecturerProfile(c *fiber.Ctx) error {
	var req struct {
		UserID     string `json:"user_id" validate:"required"`
		LecturerID string `json:"lecturer_id" validate:"required"` // NIP
		Department string `json:"department" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Validate required fields
	if req.UserID == "" || req.LecturerID == "" || req.Department == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "user_id, lecturer_id, and department are required")
	}

	// Verify user exists and has Dosen role
	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "user not found")
	}

	roleRepo := repository.NewRoleRepository()
	role, err := roleRepo.FindByID(user.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "role not found for this user")
	}

	// Check if role is Dosen (case-insensitive and flexible)
	isDosenRole := role.Name == "Dosen" || role.Name == "Lecturer" ||
		role.Name == "Dosen Wali" || role.Name == "dosen" || role.Name == "lecturer"

	if !isDosenRole {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "user must have Dosen/Lecturer role, current role: "+role.Name)
	}

	lecturer := &models.Lecturer{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		LecturerID: req.LecturerID,
		Department: req.Department,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.lecturerRepo.Create(lecturer); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create lecturer profile")
	}

	return utils.SuccessResponse(c, "lecturer profile created successfully", lecturer)
}

// UpdateLecturerProfile handles updating lecturer profile
// @Summary Update lecturer profile
// @Description Update an existing lecturer profile
// @Tags Lecturers
// @Accept json
// @Produce json
// @Param id path string true "Lecturer ID"
// @Param body body models.Lecturer true "Lecturer data"
// @Success 200 {object} models.Lecturer
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /lecturers/{id} [put]
// @Security Bearer
func (s *lecturerServiceImpl) UpdateLecturerProfile(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	var req struct {
		Department string `json:"department"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	lecturer, err := s.lecturerRepo.FindByLecturerID(lecturerID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "lecturer not found")
	}

	if req.Department != "" {
		lecturer.Department = req.Department
	}

	lecturer.UpdatedAt = time.Now()

	if err := s.lecturerRepo.Update(lecturerID, lecturer); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update lecturer profile")
	}

	return utils.SuccessResponse(c, "lecturer profile updated successfully", lecturer)
}

// VerifyAchievement handles achievement verification (Dosen Wali only)
// This endpoint is documented in achievement_service.go as /achievements/{id}/verify
func (s *lecturerServiceImpl) VerifyAchievement(c *fiber.Ctx) error {
	if c.Locals("role") != "Dosen Wali" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "only dosen wali can verify achievements")
	}

	achievementID := c.Params("id")
	dosenID := c.Locals("userID").(string)

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

// RejectAchievement handles achievement rejection (Dosen Wali only)
// This endpoint is documented in achievement_service.go as /achievements/{id}/reject
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
// @Tags Lecturers
// @Produce json
// @Success 200 {array} models.AchievementReference
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lecturer/achievements [get]
// @Security Bearer
func (s *lecturerServiceImpl) GetGuidedStudentsAchievements(c *fiber.Ctx) error {
	if c.Locals("role") != "Dosen Wali" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "only dosen wali can access this endpoint")
	}

	dosenID := c.Locals("userID").(string)

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
			"id":               ach.ID,
			"student_id":       ach.StudentID,
			"title":            mongoAch.Title,
			"description":      mongoAch.Description,
			"achievement_type": mongoAch.AchievementType,
			"details":          mongoAch.Details,
			"tags":             mongoAch.Tags,
			"points":           mongoAch.Points,
			"status":           ach.Status,
			"submitted_at":     ach.SubmittedAt,
			"verified_at":      ach.VerifiedAt,
			"verified_by":      ach.VerifiedBy,
			"rejection_note":   ach.RejectionNote,
			"created_at":       ach.CreatedAt,
		})
	}

	return utils.SuccessResponse(c, "guided students achievements retrieved", result)
}

// ListLecturers handles listing all lecturers
// @Summary List all lecturers
// @Description Get paginated list of all lecturers. Dosen/Dosen Wali only sees their own profile.
// @Tags Lecturers
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lecturers [get]
// @Security Bearer
func (s *lecturerServiceImpl) ListLecturers(c *fiber.Ctx) error {
	// Get logged in user
	userIDInterface := c.Locals("userID")
	if userIDInterface == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user not authenticated")
	}
	userID := userIDInterface.(string)
	loggedInUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get user info")
	}

	// Get user role
	role, err := s.roleRepo.FindByID(loggedInUser.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get user role")
	}

	// Fetch lecturers based on role
	var lecturers []models.Lecturer
	if role.Name == "Dosen" || role.Name == "Dosen Wali" {
		// Dosen/Dosen Wali can only see their own profile
		lecturer, err := s.lecturerRepo.FindByUserID(userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "lecturer profile not found")
		}
		lecturers = []models.Lecturer{*lecturer}
	} else {
		// Admin can see all lecturers
		lecturers, err = s.lecturerRepo.FindAll()
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to list lecturers")
		}
	}

	return utils.SuccessResponse(c, "lecturers retrieved successfully", fiber.Map{
		"data": lecturers,
	})
}

// GetAdvisees handles getting lecturer's advisees
// @Summary Get lecturer's advisees
// @Description Get list of students guided by a lecturer. Dosen/Dosen Wali can only access their own advisees.
// @Tags Lecturers
// @Produce json
// @Param id path string true "Lecturer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lecturers/{id}/advisees [get]
// @Security Bearer
func (s *lecturerServiceImpl) GetAdvisees(c *fiber.Ctx) error {
	id := c.Params("id")

	// Verify lecturer exists and get lecturer.id
	lecturer, err := s.lecturerRepo.FindByID(id)
	if err != nil {
		// Try finding by UserID
		lecturer, err = s.lecturerRepo.FindByUserID(id)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "lecturer not found")
		}
	}

	// Check ownership for Dosen/Dosen Wali
	userIDInterface := c.Locals("userID")
	if userIDInterface == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "user not authenticated")
	}
	userID := userIDInterface.(string)
	loggedInUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get user info")
	}

	role, err := s.roleRepo.FindByID(loggedInUser.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get user role")
	}

	// If Dosen/Dosen Wali, verify ownership
	if role.Name == "Dosen" || role.Name == "Dosen Wali" {
		myLecturer, err := s.lecturerRepo.FindByUserID(userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "lecturer profile not found")
		}
		if lecturer.ID != myLecturer.ID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only access your own advisees")
		}
	}

	// Get advisees (students with this lecturer as advisor) using lecturer.id
	students, errStud := s.studentRepo.FindByAdvisorID(lecturer.ID)
	if errStud != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get advisees")
	}

	return utils.SuccessResponse(c, "advisees retrieved successfully", fiber.Map{
		"data": students,
	})
}
