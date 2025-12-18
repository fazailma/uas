package service

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/utils"
)

// UserService defines all user management operations
type UserService interface {
	CreateUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	UpdateUserRole(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	ListUsers(c *fiber.Ctx) error
	GetAllAchievements(c *fiber.Ctx) error
	GetStudentAchievements(c *fiber.Ctx) error
	GetAchievementStats(c *fiber.Ctx) error
}

type userServiceImpl struct {
	userRepo             *repository.UserRepository
	roleRepo             *repository.RoleRepository
	studentRepo          *repository.StudentRepository
	achievementRepo      *repository.AchievementRepository
	mongoAchievementRepo *repository.MongoAchievementRepository
}

func NewUserService() UserService {
	return &userServiceImpl{
		userRepo:             repository.NewUserRepository(),
		roleRepo:             repository.NewRoleRepository(),
		studentRepo:          repository.NewStudentRepository(),
		achievementRepo:      repository.NewAchievementRepository(),
		mongoAchievementRepo: repository.NewMongoAchievementRepository(),
	}
}

// FunctionName godoc
// @Summary Create user
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param body body models.CreateUserRequest true "User data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Router /users [post]
// @Security Bearer
func (s *userServiceImpl) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Username == "" || req.Email == "" || req.FullName == "" || req.RoleID == "" || req.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "username, email, full_name, role_id, and password are required")
	}

	if user, _ := s.userRepo.FindByUsername(req.Username); user != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "username already exists")
	}
	if user, _ := s.userRepo.FindByEmail(req.Email); user != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "email already exists")
	}

	role, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "role not found")
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: utils.HashPassword(req.Password),
		RoleID:       req.RoleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.userRepo.Create(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create user")
	}

	// Auto-create student profile if role is Mahasiswa and student_id provided
	if role.Name == "Mahasiswa" && req.StudentID != "" {
		if err := s.studentRepo.Create(&models.Student{
			ID:           uuid.New().String(),
			UserID:       user.ID,
			StudentID:    req.StudentID,
			ProgramStudy: req.ProgramStudy,
			AcademicYear: req.AcademicYear,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create student profile")
		}
	}

	return utils.CreatedResponse(c, "user created successfully", &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     role.Name,
		RoleID:   role.ID,
		IsActive: user.IsActive,
	})
}

// FunctionName godoc
// @Summary Update user
// @Description Update user information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body models.UpdateUserRequest true "Updated user data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [put]
// @Security Bearer
func (s *userServiceImpl) UpdateUser(c *fiber.Ctx) error {
	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	user, err := s.userRepo.FindByID(c.Params("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "user not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to find user")
	}

	if req.Email != "" && req.Email != user.Email {
		if existing, _ := s.userRepo.FindByEmail(req.Email); existing != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "email already exists")
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.RoleID != "" && req.RoleID != user.RoleID {
		if _, err := s.roleRepo.FindByID(req.RoleID); err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "role not found")
		}
		user.RoleID = req.RoleID
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update user")
	}

	role, _ := s.roleRepo.FindByID(user.RoleID)
	return utils.SuccessResponse(c, "user updated successfully", &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     role.Name,
		RoleID:   role.ID,
		IsActive: user.IsActive,
	})
}

// FunctionName godoc
// @Summary Delete user
// @Description Permanently delete a user account (hard delete)
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [delete]
// @Security Bearer
func (s *userServiceImpl) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Check if user exists first
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "user not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to find user")
	}

	// Hard delete the user
	if err := s.userRepo.Delete(userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to delete user")
	}

	return utils.DeletedResponse(c, "user permanently deleted successfully")
}

// FunctionName godoc
// @Summary Get user by ID
// @Description Retrieve a specific user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
// @Security Bearer
func (s *userServiceImpl) GetUserByID(c *fiber.Ctx) error {
	user, err := s.userRepo.FindByID(c.Params("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "user not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to find user")
	}

	role, _ := s.roleRepo.FindByID(user.RoleID)
	return utils.SuccessResponse(c, "user retrieved successfully", &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     role.Name,
		RoleID:   role.ID,
		IsActive: user.IsActive,
	})
}

// FunctionName godoc
// @Summary List users
// @Description Get paginated list of all users
// @Tags Users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} models.UserResponse
// @Router /users [get]
// @Security Bearer
func (s *userServiceImpl) ListUsers(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)
	users, total, err := s.userRepo.FindAll(pagination.Page, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch users")
	}

	var responses []*models.UserResponse
	for _, user := range users {
		role, _ := s.roleRepo.FindByID(user.RoleID)
		responses = append(responses, &models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     role.Name,
			RoleID:   role.ID,
			IsActive: user.IsActive,
		})
	}

	return utils.PaginatedResponse(c, fiber.Map{"users": responses}, total, pagination.Page, pagination.Limit)
}

// This endpoint is documented in achievement_service.go as GET /achievements
func (s *userServiceImpl) GetAllAchievements(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)
	achievements, total, err := s.achievementRepo.FindAllWithPagination(pagination.Page, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to fetch achievements")
	}

	var results []*models.AchievementDetailResponse
	for i := range achievements {
		ref := &achievements[i]
		mongoAchievement, err := s.mongoAchievementRepo.FindByID(context.Background(), ref.MongoAchievementID)
		if err != nil {
			continue
		}

		results = append(results, &models.AchievementDetailResponse{
			ID:              mongoAchievement.ID.Hex(),
			Title:           mongoAchievement.Title,
			Description:     mongoAchievement.Description,
			AchievementType: mongoAchievement.AchievementType,
			Details:         mongoAchievement.Details,
			Tags:            mongoAchievement.Tags,
			Points:          mongoAchievement.Points,
			Status:          ref.Status,
			StudentID:       mongoAchievement.StudentID,
			CreatedAt:       mongoAchievement.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return utils.PaginatedResponse(c, fiber.Map{"achievements": results}, total, pagination.Page, pagination.Limit)
}

// This endpoint is documented in achievement_service.go as GET /achievements/stats
func (s *userServiceImpl) GetAchievementStats(c *fiber.Ctx) error {
	stats := make(map[string]interface{})

	totalCount, err := s.achievementRepo.CountTotal()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to count achievements")
	}
	stats["total"] = totalCount

	countByStatus, err := s.achievementRepo.CountByStatus()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to count by status")
	}

	var byStatus []map[string]interface{}
	for status, count := range countByStatus {
		byStatus = append(byStatus, map[string]interface{}{
			"status": status,
			"count":  count,
		})
	}
	stats["by_status"] = byStatus

	return utils.SuccessResponse(c, "statistics retrieved successfully", stats)
}

// FunctionName godoc
// @Summary Update user role
// @Description Update role for a user (Admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body map[string]interface{} true "Role data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id}/role [put]
// @Security Bearer
func (s *userServiceImpl) UpdateUserRole(c *fiber.Ctx) error {
	userID := c.Params("id")

	var req struct {
		RoleID string `json:"role_id" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "user not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to find user")
	}

	// Verify role exists
	role, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil || role == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid role")
	}

	// Update user role
	user.RoleID = req.RoleID
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update user role")
	}

	return utils.SuccessResponse(c, "user role updated successfully", nil)
}

func (s *userServiceImpl) GetStudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")

	achievementRepo := repository.NewAchievementRepository()
	achievements, err := achievementRepo.FindByStudentID(studentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get student achievements")
	}

	return utils.SuccessResponse(c, "achievements retrieved successfully", fiber.Map{
		"data": achievements,
	})
}
