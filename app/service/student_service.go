package service

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"UAS/app/models"
	"UAS/app/repository"
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

// FunctionName godoc
// @Summary Create student profile
// @Description Create a new student profile for a user
// @Tags Students
// @Accept json
// @Produce json
// @Param body body models.Student true "Student data"
// @Success 201 {object} models.Student
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
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

// FunctionName godoc
// @Summary Update student profile
// @Description Update an existing student profile
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param body body models.Student true "Student data"
// @Success 200 {object} models.Student
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
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

// FunctionName godoc
// @Summary Set student advisor
// @Description Assign a lecturer as advisor for a student
// @Tags Students
// @Accept json
// @Produce json
// @Param id path string true "Student UUID ID (from database, not NIM)"
// @Param body body map[string]interface{} true "Advisor lecturer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
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
	student, err := s.studentRepo.FindByID(studentUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to find student")
	}

	// Try to find lecturer by ID first, then by UserID for flexibility
	var lecturer *models.Lecturer
	lecturer, err = s.lecturerRepo.FindByID(req.AdvisorID)
	if err != nil {
		// Try finding by UserID as fallback
		lecturer, err = s.lecturerRepo.FindByUserID(req.AdvisorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ErrorResponse(c, fiber.StatusNotFound, "lecturer not found with ID: "+req.AdvisorID)
			}
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "database error finding lecturer: "+err.Error())
		}
	}

	// Verify the user associated with this lecturer has Dosen/Dosen Wali role
	advisor, err := s.userRepo.FindByID(lecturer.UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get lecturer user info")
	}

	role, err := s.roleRepo.FindByID(advisor.RoleID)
	if err != nil || role == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "role not found for lecturer")
	}

	// Check if advisor has Dosen role (flexible matching)
	isDosenRole := role.Name == "Dosen" || role.Name == "Lecturer" ||
		role.Name == "Dosen Wali" || role.Name == "dosen" || role.Name == "lecturer"

	if !isDosenRole {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "advisor must have Dosen/Lecturer role, current role: "+role.Name)
	}

	// Set advisor_id to lecturer.ID (not user_id)
	student.AdvisorID = lecturer.ID
	student.UpdatedAt = time.Now()

	if err := s.studentRepo.Update(studentUUID, student); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to set advisor")
	}

	return utils.SuccessResponse(c, "advisor set successfully", fiber.Map{
		"student_id":   student.ID,
		"advisor_id":   lecturer.ID,
		"advisor_name": advisor.FullName,
		"advisor_nip":  lecturer.LecturerID,
	})
}

// FunctionName godoc
// @Summary List all students
// @Description Get paginated list of all students. Dosen Wali only sees their advisees.
// @Tags Students
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /students [get]
// @Security Bearer
func (s *studentServiceImpl) ListStudents(c *fiber.Ctx) error {
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

	// Fetch students based on role
	var students []models.Student
	if role.Name == "Dosen Wali" {
		// Dosen Wali can only see their advisees
		lecturer, err := s.lecturerRepo.FindByUserID(userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "lecturer profile not found")
		}
		students, err = s.studentRepo.FindByAdvisorID(lecturer.ID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch advisees")
		}
	} else {
		// Admin and others can see all students
		students, err = s.studentRepo.FindAll()
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch students")
		}
	}

	// Enrich with user data
	var results []fiber.Map
	for _, student := range students {
		user, err := s.userRepo.FindByID(student.UserID)
		if err != nil {
			continue
		}

		// Get advisor name if exists
		advisorName := "-" // Default jika belum ada advisor
		if student.AdvisorID != "" && student.AdvisorID != "null" {
			advisor, err := s.lecturerRepo.FindByID(student.AdvisorID)
			if err == nil && advisor != nil {
				advisorUser, err := s.userRepo.FindByID(advisor.UserID)
				if err == nil && advisorUser != nil {
					advisorName = advisorUser.FullName
				}
			}
		}

		results = append(results, fiber.Map{
			"id":            student.ID,
			"user_id":       student.UserID,
			"student_id":    student.StudentID,
			"full_name":     user.FullName,
			"email":         user.Email,
			"program_study": student.ProgramStudy,
			"academic_year": student.AcademicYear,
			"advisor_id":    student.AdvisorID,
			"advisor_name":  advisorName,
		})
	}

	return utils.SuccessResponse(c, "students retrieved successfully", fiber.Map{
		"students": results,
		"total":    len(results),
	})
}

// FunctionName godoc
// @Summary Get student detail
// @Description Get detailed information of a student. Dosen Wali can only access their advisees.
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} models.Student
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /students/{id} [get]
// @Security Bearer
func (s *studentServiceImpl) GetStudent(c *fiber.Ctx) error {
	id := c.Params("id")

	// Try to find by ID first, then by UserID
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		// Try finding by UserID
		student, err = s.studentRepo.FindByUserID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
			}
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get student")
		}
	}

	// Check ownership for Dosen Wali
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

	// If Dosen Wali, verify ownership
	if role.Name == "Dosen Wali" {
		lecturer, err := s.lecturerRepo.FindByUserID(userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "lecturer profile not found")
		}
		if student.AdvisorID != lecturer.ID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only access your own advisees")
		}
	}

	return utils.SuccessResponse(c, "student retrieved successfully", student)
}

// FunctionName godoc
// @Summary Get student achievements
// @Description Get all achievements for a specific student. Dosen Wali can only access their advisees' achievements.
// @Tags Students
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /students/{id}/achievements [get]
// @Security Bearer
func (s *studentServiceImpl) GetStudentAchievements(c *fiber.Ctx) error {
	id := c.Params("id")

	// Verify student exists and get user_id
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		// Try finding by UserID
		student, err = s.studentRepo.FindByUserID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
			}
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to verify student")
		}
	}

	// Check ownership for Dosen Wali
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

	// If Dosen Wali, verify ownership
	if role.Name == "Dosen Wali" {
		lecturer, err := s.lecturerRepo.FindByUserID(userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "lecturer profile not found")
		}
		if student.AdvisorID != lecturer.ID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only access achievements of your own advisees")
		}
	}

	// Get achievements from achievement repository using UserID
	achievementRepo := repository.NewAchievementRepository()
	achievements, errAch := achievementRepo.FindByStudentID(student.UserID)
	if errAch != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get student achievements")
	}

	return utils.SuccessResponse(c, "achievements retrieved successfully", fiber.Map{
		"data": achievements,
	})
}
