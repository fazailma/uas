package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/database"
	"UAS/utils"

	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo     *repository.UserRepository
	studentRepo  *repository.StudentRepository
	lecturerRepo *repository.LecturerRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		studentRepo:  repository.NewStudentRepository(),
		lecturerRepo: repository.NewLecturerRepository(),
	}
}

// Login authenticates user and returns JWT token with new format
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

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(userWithPerms)
	if err != nil {
		return nil, err
	}

	// Build user profile
	userProfile := models.UserProfile{
		ID:          userWithPerms.ID,
		Username:    userWithPerms.Username,
		FullName:    userWithPerms.FullName,
		Role:        role.Name,
		Permissions: permissionNames,
	}

	// Build response in new format
	response := &models.LoginResponse{
		Status: "success",
		Data: models.LoginResponseData{
			Token:        token,
			RefreshToken: refreshToken,
			User:         userProfile,
		},
	}

	return response, nil
}

// Register creates a new user and returns user ID
func (s *AuthService) Register(reg *models.RegisterRequest) (string, error) {
	// Validate input
	if reg.Username == "" {
		return "", errors.New("username is required")
	}
	if reg.Password == "" {
		return "", errors.New("password is required")
	}

	// Check uniqueness
	if _, err := s.userRepo.FindByUsername(reg.Username); err == nil {
		return "", errors.New("username already exists")
	}

	// Hash password
	hashed := utils.HashPassword(reg.Password)

	// Create user model
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     reg.Username,
		Email:        reg.Email,
		PasswordHash: hashed,
		FullName:     reg.FullName,
		RoleID:       reg.RoleID,
		IsActive:     true,
	}

	// Save user
	if err := s.userRepo.Create(user); err != nil {
		return "", err
	}

	// Auto-create Student or Lecturer profile based on role
	// Get role dari RoleID
	var role models.Role
	if err := database.DB.Where("id = ?", user.RoleID).First(&role).Error; err != nil {
		// Role tidak ditemukan, tapi user sudah dibuat, jadi return user.ID saja
		return user.ID, nil
	}

	// Jika role adalah "Mahasiswa", buat Student record
	if role.Name == "Mahasiswa" {
		// Generate StudentID (NIM) format: TAHUN + NOMOR (e.g., 2025001)
		studentID := generateStudentID()

		// Auto-assign advisor dengan load-balancing
		advisorID := assignAdvisor()

		student := &models.Student{
			ID:           uuid.New().String(),
			UserID:       user.ID,
			StudentID:    studentID,
			ProgramStudy: "",
			AcademicYear: "",
			AdvisorID:    advisorID,
		}
		if err := s.studentRepo.Create(student); err != nil {
			// Log error tapi jangan gagalkan registration
		}
	}

	// Jika role adalah "Dosen Wali", buat Lecturer record
	if role.Name == "Dosen Wali" {
		// Generate LecturerID (NIP) format: NOMOR (e.g., 001 + timestamp)
		lecturerID := generateLecturerID()

		lecturer := &models.Lecturer{
			ID:         uuid.New().String(),
			UserID:     user.ID,
			LecturerID: lecturerID,
			Department: "",
		}
		if err := s.lecturerRepo.Create(lecturer); err != nil {
			// Log error tapi jangan gagalkan registration
		}
	}

	return user.ID, nil
}

// generateStudentID generates a unique Student ID (NIM)
// Format: YEAR + SEQUENTIAL NUMBER (e.g., 20250001)
func generateStudentID() string {
	year := time.Now().Year()

	// Count existing students untuk tahun ini
	var count int64
	database.DB.Model(&models.Student{}).
		Where("student_id LIKE ?", fmt.Sprintf("%d%%", year)).
		Count(&count)

	// Generate: YEAR + PAD NUMBER (e.g., 20250001)
	return fmt.Sprintf("%d%04d", year, count+1)
}

// generateLecturerID generates a unique Lecturer ID (NIP)
// Format: SEQUENTIAL NUMBER + TIMESTAMP (e.g., 001-1700000000)
func generateLecturerID() string {
	// Count existing lecturers
	var count int64
	database.DB.Model(&models.Lecturer{}).Count(&count)

	// Generate: PAD NUMBER + TIMESTAMP (e.g., 001-1700000000)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%03d-%d", count+1, timestamp)
}

// assignAdvisor assigns a student to a lecturer with load-balancing
// Selects the lecturer with the fewest students
// If no lecturers exist, returns empty string
func assignAdvisor() string {
	// Get all lecturers
	var lecturers []models.Lecturer
	if err := database.DB.Find(&lecturers).Error; err != nil {
		return "" // Return empty if query fails
	}

	if len(lecturers) == 0 {
		return "" // No lecturers available
	}

	// For each lecturer, count how many students they advise
	var selectedLecturer models.Lecturer
	var minStudentCount int64 = 999999 // Large number to start

	for _, lecturer := range lecturers {
		var count int64
		database.DB.Model(&models.Student{}).
			Where("advisor_id = ?", lecturer.ID).
			Count(&count)

		// Select lecturer with fewest students
		if count < minStudentCount {
			minStudentCount = count
			selectedLecturer = lecturer
		}
	}

	return selectedLecturer.ID
}
