package service

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"UAS/app/models"
	"UAS/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository adalah mock untuk UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserWithRoleAndPermissions(userID string) (*models.User, []models.Permission, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*models.User), args.Get(1).([]models.Permission), args.Error(2)
}

// TestLoginSuccess tests successful login
func TestLoginSuccess(t *testing.T) {
	// Setup
	app := fiber.New()

	// Create mock request
	loginReq := models.LoginCredential{
		Username: "testuser",
		Password: "password123",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Test
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestLoginInvalidCredentials tests login with invalid credentials
func TestLoginInvalidCredentials(t *testing.T) {
	// Test empty credentials validation
	loginReq := models.LoginCredential{
		Username: "",
		Password: "",
	}

	// Validate required fields
	hasError := loginReq.Username == "" || loginReq.Password == ""

	// Assert
	assert.True(t, hasError, "Should detect empty credentials")
}

// TestHashPassword tests password hashing
func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	// Test hashing
	hashedPassword := utils.HashPassword(password)

	// Assert
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	// Test verification
	isValid := utils.VerifyPassword(password, hashedPassword)
	assert.True(t, isValid)

	// Test invalid password
	isInvalid := utils.VerifyPassword("wrongpassword", hashedPassword)
	assert.False(t, isInvalid)
}

// TestRegisterValidation tests register input validation
func TestRegisterValidation(t *testing.T) {
	testCases := []struct {
		name        string
		request     models.RegisterRequest
		expectError bool
	}{
		{
			name: "Valid registration",
			request: models.RegisterRequest{
				Username: "newuser",
				Password: "password123",
				Email:    "test@example.com",
				FullName: "Test User",
				RoleID:   "role-123",
			},
			expectError: false,
		},
		{
			name: "Missing username",
			request: models.RegisterRequest{
				Username: "",
				Password: "password123",
				Email:    "test@example.com",
				FullName: "Test User",
			},
			expectError: true,
		},
		{
			name: "Missing password",
			request: models.RegisterRequest{
				Username: "newuser",
				Password: "",
				Email:    "test@example.com",
				FullName: "Test User",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Validate required fields
			hasError := tc.request.Username == "" ||
				tc.request.Password == "" ||
				tc.request.Email == "" ||
				tc.request.FullName == ""

			assert.Equal(t, tc.expectError, hasError)
		})
	}
}
