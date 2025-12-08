package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/database"
	"UAS/helpers"
	"UAS/utils"
)

// AdminService handles admin operations like user management and achievement viewing
type AdminService struct {
	userRepo             *repository.UserRepository
	roleRepo             *repository.RoleRepository
	studentRepo          *repository.StudentRepository
	lecturerRepo         *repository.LecturerRepository
	achievementRepo      *repository.AchievementRepository
	mongoAchievementRepo *repository.MongoAchievementRepository
}

// NewAdminService creates a new instance of AdminService
func NewAdminService() *AdminService {
	return &AdminService{
		userRepo:             repository.NewUserRepository(),
		roleRepo:             repository.NewRoleRepository(),
		studentRepo:          repository.NewStudentRepository(),
		lecturerRepo:         repository.NewLecturerRepository(),
		achievementRepo:      repository.NewAchievementRepository(),
		mongoAchievementRepo: repository.NewMongoAchievementRepository(),
	}
}

// ===== FR-009: Manage Users - Logic Methods (Private) =====

// createUser is the logic method that creates a new user with role assignment
func (s *AdminService) createUser(req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Validate input
	if req.Username == "" || req.Email == "" || req.FullName == "" || req.RoleID == "" {
		return nil, errors.New("username, email, full_name, and role_id are required")
	}

	// Check if username already exists
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingEmail, _ := s.userRepo.FindByEmail(req.Email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Verify role exists
	role, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Hash password
	if req.Password == "" {
		return nil, errors.New("password is required")
	}
	hashedPassword := utils.HashPassword(req.Password)

	// Create user
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: hashedPassword,
		RoleID:       req.RoleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// If user type is student, create student profile
	if req.StudentID != "" {
		student := &models.Student{
			ID:           uuid.New().String(),
			UserID:       user.ID,
			StudentID:    req.StudentID,
			ProgramStudy: req.ProgramStudy,
			AcademicYear: req.AcademicYear,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := s.studentRepo.Create(student); err != nil {
			return nil, fmt.Errorf("failed to create student profile: %w", err)
		}
	}

	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     role.Name,
		RoleID:   role.ID,
		IsActive: user.IsActive,
	}, nil
}

// updateUser updates user information
func (s *AdminService) updateUser(userID string, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Find user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Update fields if provided
	if req.Email != "" && req.Email != user.Email {
		// Check if new email already exists
		existingEmail, _ := s.userRepo.FindByEmail(req.Email)
		if existingEmail != nil {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.RoleID != "" && req.RoleID != user.RoleID {
		// Verify role exists
		role, err := s.roleRepo.FindByID(req.RoleID)
		if err != nil {
			return nil, errors.New("role not found")
		}
		user.RoleID = role.ID
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	role, _ := s.roleRepo.FindByID(user.RoleID)

	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     role.Name,
		RoleID:   role.ID,
		IsActive: user.IsActive,
	}, nil
}

// deleteUser deletes a user
func (s *AdminService) deleteUser(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	user.IsActive = false
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// getUserByID retrieves user information by ID
func (s *AdminService) getUserByID(userID string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	role, _ := s.roleRepo.FindByID(user.RoleID)

	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     role.Name,
		RoleID:   role.ID,
		IsActive: user.IsActive,
	}, nil
}

// listUsers retrieves all users with pagination
func (s *AdminService) listUsers(page, pageSize int) ([]*models.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.userRepo.FindAll(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch users: %w", err)
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

	return responses, total, nil
}

// setAdvisor assigns an advisor (lecturer) to a student
func (s *AdminService) setAdvisor(studentID, advisorID string) error {
	// Verify student exists by NIM (student_id)
	student, err := s.studentRepo.FindByStudentID(studentID)
	if err != nil {
		return errors.New("student not found")
	}

	// Verify advisor exists and is a lecturer
	advisor, err := s.userRepo.FindByID(advisorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("advisor not found")
		}
		return fmt.Errorf("failed to find advisor: %w", err)
	}

	// Check if advisor has lecturer role
	role, _ := s.roleRepo.FindByID(advisor.RoleID)
	if role == nil || role.Name != "Dosen" {
		return errors.New("advisor must have lecturer role")
	}

	// Update student's advisor
	student.AdvisorID = advisorID
	student.UpdatedAt = time.Now()

	if err := s.studentRepo.Update(studentID, student); err != nil {
		return fmt.Errorf("failed to set advisor: %w", err)
	}

	return nil
}

// ===== FR-010: View All Achievements - Logic Methods (Private) =====

// getAllAchievements retrieves all achievements with pagination
func (s *AdminService) getAllAchievements(page, pageSize int) (interface{}, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get achievement references from PostgreSQL
	query := database.DB

	var achievements []*models.AchievementReference
	var total int64

	if err := query.Model(&models.AchievementReference{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count achievements: %w", err)
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&achievements).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch achievements: %w", err)
	}

	// Build response with MongoDB details
	var results []*models.AchievementDetailResponse
	for _, ref := range achievements {
		// Fetch from MongoDB
		mongoAchievement, err := s.mongoAchievementRepo.FindByID(context.Background(), ref.MongoAchievementID)
		if err != nil {
			continue
		}

		results = append(results, &models.AchievementDetailResponse{
			ID:          mongoAchievement.ID.Hex(),
			Title:       mongoAchievement.Title,
			Description: mongoAchievement.Description,
			Category:    mongoAchievement.Category,
			Date:        mongoAchievement.Date,
			ProofURL:    mongoAchievement.ProofURL,
			Status:      ref.Status,
			StudentID:   mongoAchievement.StudentID,
			CreatedAt:   mongoAchievement.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return results, total, nil
}

// getAchievementStats returns statistics about achievements
func (s *AdminService) getAchievementStats() (map[string]interface{}, error) {
	var stats = make(map[string]interface{})

	// Total achievements
	var totalCount int64
	if err := database.DB.Model(&models.AchievementReference{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total"] = totalCount

	// Count by status
	var statusStats []map[string]interface{}
	type StatusCount struct {
		Status string
		Count  int64
	}
	var results []StatusCount

	if err := database.DB.Model(&models.AchievementReference{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	for _, sc := range results {
		statusStats = append(statusStats, map[string]interface{}{
			"status": sc.Status,
			"count":  sc.Count,
		})
	}
	stats["by_status"] = statusStats

	return stats, nil
}

// ===== HTTP Handlers (Public) =====

// ListUsers returns all users with pagination
func (s *AdminService) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	users, total, err := s.listUsers(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.BuildErrorResponse(500, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, fiber.Map{
		"data":       users,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	}))
}

// GetUserByID returns a user by ID
func (s *AdminService) GetUserByID(c *fiber.Ctx) error {
	user, err := s.getUserByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(helpers.BuildErrorResponse(404, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, user))
}

// CreateUser creates a new user
func (s *AdminService) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "invalid request body"))
	}
	user, err := s.createUser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
	}
	return c.Status(fiber.StatusCreated).JSON(helpers.BuildCreatedResponse("user created successfully", user))
}

// UpdateUser updates a user
func (s *AdminService) UpdateUser(c *fiber.Ctx) error {
	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "invalid request body"))
	}
	user, err := s.updateUser(c.Params("id"), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("user updated successfully", user))
}

// DeleteUser deletes a user
func (s *AdminService) DeleteUser(c *fiber.Ctx) error {
	if err := s.deleteUser(c.Params("id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(helpers.BuildDeletedResponse("user deleted successfully"))
}

// SetAdvisor assigns an advisor to a student
func (s *AdminService) SetAdvisor(c *fiber.Ctx) error {
	if err := s.setAdvisor(c.Params("student_id"), c.Params("advisor_id")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("advisor set successfully", nil))
}

// GetAchievements returns all achievements
func (s *AdminService) GetAchievements(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	achievements, total, err := s.getAllAchievements(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.BuildErrorResponse(500, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, fiber.Map{
		"data":       achievements,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	}))
}

// GetAchievementStats returns achievement statistics
func (s *AdminService) GetAchievementStats(c *fiber.Ctx) error {
	stats, err := s.getAchievementStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.BuildErrorResponse(500, err.Error()))
	}
	return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, stats))
}
