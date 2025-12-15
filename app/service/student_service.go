package service

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/database"
	"UAS/utils"
)

// StudentService defines all student-related operations
type StudentService interface {
	CreateStudentProfile(c *fiber.Ctx) error
	ListStudents(c *fiber.Ctx) error
	GetStudent(c *fiber.Ctx) error
	UpdateStudentProfile(c *fiber.Ctx) error
	GetStudentAchievements(c *fiber.Ctx) error
	SetAdvisor(c *fiber.Ctx) error
}

type studentServiceImpl struct {
	studentRepo  *repository.StudentRepository
	userRepo     *repository.UserRepository
	roleRepo     *repository.RoleRepository
	lecturerRepo *repository.LecturerRepository
}

func NewStudentService() StudentService {
	return &studentServiceImpl{
		studentRepo:  repository.NewStudentRepository(),
		userRepo:     repository.NewUserRepository(),
		roleRepo:     repository.NewRoleRepository(),
		lecturerRepo: repository.NewLecturerRepository(),
	}
}

// CreateStudentProfile handles creating student profile
// @Summary Create student profile
// @Description Create a new student profile for a user
// @Tags Student
// @Accept json
// @Produce json
// @Param body body models.Student true "Student data"
// @Success 201 {object} models.Student
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students [post]
// @Security Bearer
func (s *studentServiceImpl) CreateStudentProfile(c *fiber.Ctx) error {
	var req struct {
		UserID       string `json:"user_id" validate:"required"`
		StudentID    string `json:"student_id" validate:"required"` // NIM
		ProgramStudy string `json:"program_study" validate:"required"`
		AcademicYear string `json:"academic_year" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Validate required fields
	if req.UserID == "" || req.StudentID == "" || req.ProgramStudy == "" || req.AcademicYear == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "user_id, student_id, program_study, and academic_year are required")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "user not found")
	}

	// Verify user has Mahasiswa role
	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "role not found for this user")
	}

	// Check if role is Mahasiswa (case-insensitive)
	isMahasiswaRole := role.Name == "Mahasiswa" || role.Name == "mahasiswa" || role.Name == "Student"

	if !isMahasiswaRole {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "user must have Mahasiswa role, current role: "+role.Name)
	}

	student := &models.Student{
		ID:           uuid.New().String(),
		UserID:       req.UserID,
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.studentRepo.Create(student); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create student profile")
	}

	return utils.SuccessResponse(c, "student profile created successfully", student)
}

// UpdateStudentProfile handles updating student profile
// @Summary Update student profile
// @Description Update an existing student profile
// @Tags Student
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param body body models.Student true "Student data"
// @Success 200 {object} models.Student
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /students/{id} [put]
// @Security Bearer
func (s *studentServiceImpl) UpdateStudentProfile(c *fiber.Ctx) error {
	studentID := c.Params("id")

	var req struct {
		ProgramStudy string `json:"program_study"`
		AcademicYear string `json:"academic_year"`
		AdvisorID    string `json:"advisor_id,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	student, err := s.studentRepo.FindByStudentID(studentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
	}

	if req.ProgramStudy != "" {
		student.ProgramStudy = req.ProgramStudy
	}
	if req.AcademicYear != "" {
		student.AcademicYear = req.AcademicYear
	}
	if req.AdvisorID != "" {
		student.AdvisorID = req.AdvisorID
	}

	student.UpdatedAt = time.Now()

	if err := s.studentRepo.Update(studentID, student); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update student profile")
	}

	return utils.SuccessResponse(c, "student profile updated successfully", student)
}

// SetAdvisor handles setting student advisor
// @Summary Set student advisor
// @Description Assign a lecturer as advisor for a student
// @Tags Student
// @Accept json
// @Produce json
// @Param id path string true "Student UUID ID (from database, not NIM)"
// @Param body body models.MessageResponse true "Advisor user ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /students/{id}/advisor [put]
// @Security Bearer
func (s *studentServiceImpl) SetAdvisor(c *fiber.Ctx) error {
	studentUUID := c.Params("id")

	var req struct {
		AdvisorID string `json:"advisor_id" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.AdvisorID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "advisor_id is required")
	}

	// Find student by UUID (primary key)
	var student *models.Student
	err := database.DB.Where("id = ?", studentUUID).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to find student")
	}

	advisor, err := s.userRepo.FindByID(req.AdvisorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "advisor user ID not found in database: "+req.AdvisorID+". Make sure user Dosen was created correctly with POST /api/v1/users")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "database error finding advisor: "+err.Error())
	}

	// Verify advisor has a RoleID
	if advisor.RoleID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "advisor user found but has no role assigned")
	}

	role, err := s.roleRepo.FindByID(advisor.RoleID)
	if err != nil || role == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "role not found for advisor")
	}

	// Check if advisor has Dosen role (flexible matching)
	isDosenRole := role.Name == "Dosen" || role.Name == "Lecturer" ||
		role.Name == "Dosen Wali" || role.Name == "dosen" || role.Name == "lecturer"

	if !isDosenRole {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "advisor must have Dosen/Lecturer role, current role: "+role.Name)
	}

	student.AdvisorID = req.AdvisorID
	student.UpdatedAt = time.Now()

	if err := s.studentRepo.Update(studentUUID, student); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to set advisor")
	}

	return utils.SuccessResponse(c, "advisor set successfully", nil)
}

// ListStudents handles listing all students
// @Summary List all students
// @Description Get paginated list of all students
// @Tags Students
// @Produce json
// @Success 200 {object} models.StudentListResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students [get]
// @Security Bearer
func (s *studentServiceImpl) ListStudents(c *fiber.Ctx) error {
	// For now, get all students without pagination
	// TODO: Implement pagination in StudentRepository
	return utils.ErrorResponse(c, fiber.StatusNotImplemented, "list students not yet implemented")
}

// GetStudent handles getting student detail
// @Summary Get student detail
// @Description Get detailed information of a student
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} models.Student
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students/{id} [get]
// @Security Bearer
func (s *studentServiceImpl) GetStudent(c *fiber.Ctx) error {
	studentID := c.Params("id")

	student, err := s.studentRepo.FindByUserID(studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get student")
	}

	return utils.SuccessResponse(c, "student retrieved successfully", student)
}

// GetStudentAchievements handles getting student achievements
// @Summary Get student achievements
// @Description Get all achievements for a specific student
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} models.AchievementListResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /students/{id}/achievements [get]
// @Security Bearer
func (s *studentServiceImpl) GetStudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// Verify student exists
	_, err := s.studentRepo.FindByUserID(studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to verify student")
	}

	// Get achievements from achievement repository
	achievementRepo := repository.NewAchievementRepository()
	achievements, err := achievementRepo.FindByStudentID(studentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get student achievements")
	}

	return utils.SuccessResponse(c, "achievements retrieved successfully", fiber.Map{
		"data": achievements,
	})
}
