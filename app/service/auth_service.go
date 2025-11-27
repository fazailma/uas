package service

import (
	"errors"

	"github.com/google/uuid"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/database"
	"UAS/utils"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Login authenticates user and returns JWT token
func (s *AuthService) Login(loginReq *models.LoginCredential) (*models.LoginResponse, error) {
	// Validate input
	if loginReq.Username == "" {
		return nil, errors.New("username is required")
	}
	if loginReq.Password == "" {
		return nil, errors.New("password is required")
	}

	// Find user by username
	user, err := s.userRepo.FindByUsername(loginReq.Username)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	// Verify password
	if !utils.VerifyPassword(loginReq.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Get user with role and permissions
	userWithPerms, permissions, err := s.userRepo.GetUserWithRoleAndPermissions(user.ID)
	if err != nil {
		return nil, err
	}

	// Get role - guard against empty RoleID to avoid UUID errors
	var role models.Role
	if userWithPerms.RoleID != "" {
		if err := database.DB.Where("id = ?", userWithPerms.RoleID).First(&role).Error; err != nil {
			// Log but don't fail if role not found
			role.Name = ""
		}
	}

	// Generate JWT token
	permissionNames := make([]string, len(permissions))
	for i, p := range permissions {
		permissionNames[i] = p.Name
	}

	token, err := utils.GenerateJWT(userWithPerms, role, permissionNames)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &models.LoginResponse{
		Token:   token,
		UserID:  userWithPerms.ID,
		Message: "Login successful",
	}

	return response, nil
}

// Register creates a new user and returns JWT token (auto-login)
func (s *AuthService) Register(reg *models.LoginCredential) (*models.LoginResponse, error) {
	// Validate input
	if reg.Username == "" {
		return nil, errors.New("username is required")
	}
	if reg.Password == "" {
		return nil, errors.New("password is required")
	}

	// Check uniqueness
	if _, err := s.userRepo.FindByUsername(reg.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashed := utils.HashPassword(reg.Password)

	// Create user model
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     reg.Username,
		Email:        "",
		PasswordHash: hashed,
		FullName:     "",
		RoleID:       "",
		IsActive:     true,
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Get user with permissions
	userWithPerms, permissions, err := s.userRepo.GetUserWithRoleAndPermissions(user.ID)
	if err != nil {
		return nil, err
	}

	// Get role if any - guard against empty RoleID
	var role models.Role
	if userWithPerms.RoleID != "" {
		if err := database.DB.Where("id = ?", userWithPerms.RoleID).First(&role).Error; err != nil {
			// Log but don't fail if role not found
			role.Name = ""
		}
	}

	// Generate JWT token
	permissionNames := make([]string, len(permissions))
	for i, p := range permissions {
		permissionNames[i] = p.Name
	}

	token, err := utils.GenerateJWT(userWithPerms, role, permissionNames)
	if err != nil {
		return nil, err
	}

	response := &models.LoginResponse{
		Token:   token,
		UserID:  userWithPerms.ID,
		Message: "Registration successful",
	}

	return response, nil
}
