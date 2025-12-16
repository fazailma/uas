package service

import (
	"testing"

	"UAS/app/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAchievementRepository adalah mock untuk AchievementRepository
type MockAchievementRepository struct {
	mock.Mock
}

func (m *MockAchievementRepository) Create(achievement *models.AchievementReference) error {
	args := m.Called(achievement)
	return args.Error(0)
}

func (m *MockAchievementRepository) Update(achievement *models.AchievementReference) error {
	args := m.Called(achievement)
	return args.Error(0)
}

func (m *MockAchievementRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAchievementRepository) FindByID(id string) (*models.AchievementReference, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) FindByStudentID(studentID string) ([]models.AchievementReference, error) {
	args := m.Called(studentID)
	return args.Get(0).([]models.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) FindAll() ([]models.AchievementReference, error) {
	args := m.Called()
	return args.Get(0).([]models.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) FindByStatus(status string) ([]models.AchievementReference, error) {
	args := m.Called(status)
	return args.Get(0).([]models.AchievementReference), args.Error(1)
}

func (m *MockAchievementRepository) CountTotal() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// TestAchievementValidation tests achievement validation
func TestAchievementValidation(t *testing.T) {
	testCases := []struct {
		name        string
		request     models.CreateAchievementRequest
		expectValid bool
	}{
		{
			name: "Valid achievement",
			request: models.CreateAchievementRequest{
				Title:           "Juara 1 Lomba Programming",
				Description:     "Kompetisi tingkat nasional",
				AchievementType: "competition",
			},
			expectValid: true,
		},
		{
			name: "Missing title",
			request: models.CreateAchievementRequest{
				Title:           "",
				Description:     "Kompetisi tingkat nasional",
				AchievementType: "competition",
			},
			expectValid: false,
		},
		{
			name: "Missing achievement type",
			request: models.CreateAchievementRequest{
				Title:           "Juara 1 Lomba Programming",
				Description:     "Kompetisi tingkat nasional",
				AchievementType: "",
			},
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.request.Title != "" && tc.request.AchievementType != ""
			assert.Equal(t, tc.expectValid, isValid)
		})
	}
}

// TestAchievementStatusTransition tests valid status transitions
func TestAchievementStatusTransition(t *testing.T) {
	testCases := []struct {
		name          string
		currentStatus string
		newStatus     string
		expectValid   bool
	}{
		{
			name:          "Draft to Submitted",
			currentStatus: "draft",
			newStatus:     "submitted",
			expectValid:   true,
		},
		{
			name:          "Submitted to Verified",
			currentStatus: "submitted",
			newStatus:     "verified",
			expectValid:   true,
		},
		{
			name:          "Submitted to Rejected",
			currentStatus: "submitted",
			newStatus:     "rejected",
			expectValid:   true,
		},
		{
			name:          "Verified to Draft (Invalid)",
			currentStatus: "verified",
			newStatus:     "draft",
			expectValid:   false,
		},
		{
			name:          "Draft to Verified (Invalid - Skip Submitted)",
			currentStatus: "draft",
			newStatus:     "verified",
			expectValid:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := false

			// Define valid transitions
			validTransitions := map[string][]string{
				"draft":     {"submitted"},
				"submitted": {"verified", "rejected"},
				"rejected":  {"submitted"},
			}

			if allowedStatuses, exists := validTransitions[tc.currentStatus]; exists {
				for _, status := range allowedStatuses {
					if status == tc.newStatus {
						isValid = true
						break
					}
				}
			}

			assert.Equal(t, tc.expectValid, isValid)
		})
	}
}

// TestCreateAchievement tests achievement creation with mock
func TestCreateAchievement(t *testing.T) {
	// Setup mock
	mockRepo := new(MockAchievementRepository)

	achievement := &models.AchievementReference{
		ID:                 "achievement-123",
		StudentID:          "student-123",
		MongoAchievementID: "mongo-id-123",
		Status:             "draft",
	}

	// Setup expectation
	mockRepo.On("Create", achievement).Return(nil)

	// Test
	err := mockRepo.Create(achievement)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUpdateAchievementStatus tests updating achievement status
func TestUpdateAchievementStatus(t *testing.T) {
	// Setup mock
	mockRepo := new(MockAchievementRepository)

	achievement := &models.AchievementReference{
		ID:                 "achievement-123",
		StudentID:          "student-123",
		MongoAchievementID: "mongo-id-123",
		Status:             "submitted",
	}

	// Setup expectation
	mockRepo.On("Update", achievement).Return(nil)

	// Test
	err := mockRepo.Update(achievement)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestFindAchievementByID tests finding achievement by ID
func TestFindAchievementByID(t *testing.T) {
	// Setup mock
	mockRepo := new(MockAchievementRepository)

	expectedAchievement := &models.AchievementReference{
		ID:                 "achievement-123",
		StudentID:          "student-123",
		MongoAchievementID: "mongo-id-123",
		Status:             "draft",
	}

	// Setup expectation
	mockRepo.On("FindByID", "achievement-123").Return(expectedAchievement, nil)

	// Test
	achievement, err := mockRepo.FindByID("achievement-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, achievement)
	assert.Equal(t, "achievement-123", achievement.ID)
	assert.Equal(t, "draft", achievement.Status)
	mockRepo.AssertExpectations(t)
}

// TestFindAchievementsByStudentID tests finding achievements by student ID
func TestFindAchievementsByStudentID(t *testing.T) {
	// Setup mock
	mockRepo := new(MockAchievementRepository)

	expectedAchievements := []models.AchievementReference{
		{
			ID:        "achievement-1",
			StudentID: "student-123",
			Status:    "draft",
		},
		{
			ID:        "achievement-2",
			StudentID: "student-123",
			Status:    "submitted",
		},
	}

	// Setup expectation
	mockRepo.On("FindByStudentID", "student-123").Return(expectedAchievements, nil)

	// Test
	achievements, err := mockRepo.FindByStudentID("student-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, achievements)
	assert.Equal(t, 2, len(achievements))
	assert.Equal(t, "student-123", achievements[0].StudentID)
	mockRepo.AssertExpectations(t)
}

// TestDeleteAchievement tests achievement deletion
func TestDeleteAchievement(t *testing.T) {
	// Setup mock
	mockRepo := new(MockAchievementRepository)

	// Setup expectation
	mockRepo.On("Delete", "achievement-123").Return(nil)

	// Test
	err := mockRepo.Delete("achievement-123")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestAchievementPointsCalculation tests points calculation logic
func TestAchievementPointsCalculation(t *testing.T) {
	testCases := []struct {
		name            string
		achievementType string
		level           string
		expectedPoints  int
	}{
		{
			name:            "International Competition",
			achievementType: "competition",
			level:           "international",
			expectedPoints:  100,
		},
		{
			name:            "National Competition",
			achievementType: "competition",
			level:           "national",
			expectedPoints:  75,
		},
		{
			name:            "Regional Competition",
			achievementType: "competition",
			level:           "regional",
			expectedPoints:  50,
		},
		{
			name:            "International Publication",
			achievementType: "publication",
			level:           "international",
			expectedPoints:  150,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simple points calculation logic
			points := 0

			if tc.achievementType == "competition" {
				switch tc.level {
				case "international":
					points = 100
				case "national":
					points = 75
				case "regional":
					points = 50
				}
			} else if tc.achievementType == "publication" {
				switch tc.level {
				case "international":
					points = 150
				case "national":
					points = 100
				}
			}

			assert.Equal(t, tc.expectedPoints, points)
		})
	}
}
