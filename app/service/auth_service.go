package service

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/database"
	"UAS/utils"
)

// AuthService defines all authentication operations
type AuthService interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
}

type authServiceImpl struct {
	userRepo     *repository.UserRepository
	studentRepo  *repository.StudentRepository
	lecturerRepo *repository.LecturerRepository
}

func NewAuthService() AuthService {
	return &authServiceImpl{
		userRepo:     repository.NewUserRepository(),
		studentRepo:  repository.NewStudentRepository(),
		lecturerRepo: repository.NewLecturerRepository(),
	}
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginCredential true \"Login credentials\"
// @Success 200 {object} fiber.Map \"token and user data\"
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (s *authServiceImpl) Login(c *fiber.Ctx) error {
	var req models.LoginCredential
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Username == "" || req.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "username and password are required")
	}

	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil || !user.IsActive {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	if !utils.VerifyPassword(req.Password, user.PasswordHash) {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	userWithPerms, permissions, err := s.userRepo.GetUserWithRoleAndPermissions(user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get user permissions")
	}

	var role models.Role
	if userWithPerms.RoleID != "" {
		if err := database.DB.Where("id = ?", userWithPerms.RoleID).First(&role).Error; err != nil {
			role.Name = ""
		}
	}

	permissionNames := make([]string, len(permissions))
	for i, p := range permissions {
		permissionNames[i] = p.Name
	}

	token, err := utils.GenerateJWT(userWithPerms, role, permissionNames)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to generate token")
	}

	refreshToken, err := utils.GenerateRefreshToken(userWithPerms)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to generate refresh token")
	}

	response := &models.LoginResponse{
		Status: "success",
		Data: models.LoginResponseData{
			Token:        token,
			RefreshToken: refreshToken,
			User: models.UserProfile{
				ID:          userWithPerms.ID,
				Username:    userWithPerms.Username,
				FullName:    userWithPerms.FullName,
				Role:        role.Name,
				Permissions: permissionNames,
			},
		},
	}

	return utils.SuccessResponse(c, "login successful", response)
}

// Register handler
func (s *authServiceImpl) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Username == "" || req.Password == "" || req.Email == "" || req.FullName == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "username, password, email, and full_name are required")
	}

	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "username already exists")
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: utils.HashPassword(req.Password),
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to register user")
	}

	var role models.Role
	if err := database.DB.Where("id = ?", user.RoleID).First(&role).Error; err != nil {
		return utils.CreatedResponse(c, "user registered successfully", fiber.Map{"user_id": user.ID})
	}

	if role.Name == "Mahasiswa" {
		studentID := s.generateStudentID()
		advisorID := s.assignAdvisor()

		if err := s.studentRepo.Create(&models.Student{
			ID:           uuid.New().String(),
			UserID:       user.ID,
			StudentID:    studentID,
			ProgramStudy: "",
			AcademicYear: "",
			AdvisorID:    advisorID,
		}); err != nil {
			// Log error tapi jangan gagalkan registration
		}
	}

	if role.Name == "Dosen Wali" {
		lecturerID := s.generateLecturerID()

		if err := s.lecturerRepo.Create(&models.Lecturer{
			ID:         uuid.New().String(),
			UserID:     user.ID,
			LecturerID: lecturerID,
			Department: "",
		}); err != nil {
			// Log error tapi jangan gagalkan registration
		}
	}

	return utils.CreatedResponse(c, "user registered successfully", fiber.Map{"user_id": user.ID})
}

// Logout handler
func (s *authServiceImpl) Logout(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, "logout successful", nil)
}

// RefreshToken handler
func (s *authServiceImpl) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Token == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "token is required")
	}

	return utils.SuccessResponse(c, "token refreshed", fiber.Map{"token": "new-token-here"})
}

// GetProfile handler
func (s *authServiceImpl) GetProfile(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, "profile retrieved", fiber.Map{
		"user_id":     c.Locals("user_id"),
		"username":    c.Locals("username"),
		"email":       c.Locals("email"),
		"role":        c.Locals("role"),
		"permissions": c.Locals("permissions"),
	})
}

// generateStudentID generates a unique Student ID
func (s *authServiceImpl) generateStudentID() string {
	year := time.Now().Year()
	count, err := s.studentRepo.CountByYear(year)
	if err != nil {
		count = 0
	}
	return fmt.Sprintf("%d%04d", year, count+1)
}

// generateLecturerID generates a unique Lecturer ID
func (s *authServiceImpl) generateLecturerID() string {
	count, err := s.lecturerRepo.CountTotal()
	if err != nil {
		count = 0
	}
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%03d-%d", count+1, timestamp)
}

// assignAdvisor assigns student to lecturer with load-balancing
func (s *authServiceImpl) assignAdvisor() string {
	lecturers, err := s.lecturerRepo.FindAll()
	if err != nil || len(lecturers) == 0 {
		return ""
	}

	var selectedLecturer models.Lecturer
	var minStudentCount int64 = 999999

	for _, lecturer := range lecturers {
		count, err := s.studentRepo.CountByAdvisorID(lecturer.ID)
		if err != nil {
			continue
		}

		if count < minStudentCount {
			minStudentCount = count
			selectedLecturer = lecturer
		}
	}

	return selectedLecturer.ID
}
