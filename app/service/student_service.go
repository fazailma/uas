package service

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"UAS/app/repository"
	"UAS/utils"
)

// StudentService defines all student-related operations
type StudentService interface {
	ListStudents(c *fiber.Ctx) error
	GetStudent(c *fiber.Ctx) error
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

// SetAdvisor handles setting student advisor
// @Summary Set student advisor
// @Description Assign a lecturer as advisor for a student
// @Tags Student
// @Accept json
// @Produce json
// @Param body body map[string]string true \"Advisor data\"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /student/advisor [post]
// @Security Bearer
func (s *studentServiceImpl) SetAdvisor(c *fiber.Ctx) error {
	var req struct {
		StudentID string `json:"student_id" validate:"required"`
		AdvisorID string `json:"advisor_id" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	student, err := s.studentRepo.FindByStudentID(req.StudentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
	}

	advisor, err := s.userRepo.FindByID(req.AdvisorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "advisor not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to find advisor")
	}

	role, _ := s.roleRepo.FindByID(advisor.RoleID)
	if role == nil || role.Name != "Dosen" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "advisor must have lecturer role")
	}

	student.AdvisorID = req.AdvisorID
	student.UpdatedAt = time.Now()

	if err := s.studentRepo.Update(req.StudentID, student); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to set advisor")
	}

	return utils.SuccessResponse(c, "advisor set successfully", nil)
}

// ListStudents handles listing all students
// @Summary List all students
// @Description Get paginated list of all students
// @Tags Students
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
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
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
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
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
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
