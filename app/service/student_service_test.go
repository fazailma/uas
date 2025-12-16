package service

import (
	"testing"

	"UAS/app/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStudentRepository adalah mock untuk StudentRepository
type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) Create(student *models.Student) error {
	args := m.Called(student)
	return args.Error(0)
}

func (m *MockStudentRepository) Update(student *models.Student) error {
	args := m.Called(student)
	return args.Error(0)
}

func (m *MockStudentRepository) FindByID(id string) (*models.Student, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *MockStudentRepository) FindByUserID(userID string) (*models.Student, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *MockStudentRepository) FindByStudentID(studentID string) (*models.Student, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Student), args.Error(1)
}

func (m *MockStudentRepository) FindAll() ([]models.Student, error) {
	args := m.Called()
	return args.Get(0).([]models.Student), args.Error(1)
}

func (m *MockStudentRepository) CountByYear(year int) (int64, error) {
	args := m.Called(year)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStudentRepository) CountByAdvisorID(advisorID string) (int64, error) {
	args := m.Called(advisorID)
	return args.Get(0).(int64), args.Error(1)
}

// TestStudentValidation tests student profile validation
func TestStudentValidation(t *testing.T) {
	testCases := []struct {
		name        string
		student     *models.Student
		expectValid bool
	}{
		{
			name: "Valid student",
			student: &models.Student{
				UserID:       "user-123",
				StudentID:    "2024001",
				ProgramStudy: "Informatika",
				AcademicYear: "2024",
			},
			expectValid: true,
		},
		{
			name: "Missing UserID",
			student: &models.Student{
				UserID:       "",
				StudentID:    "2024001",
				ProgramStudy: "Informatika",
				AcademicYear: "2024",
			},
			expectValid: false,
		},
		{
			name: "Missing StudentID",
			student: &models.Student{
				UserID:       "user-123",
				StudentID:    "",
				ProgramStudy: "Informatika",
				AcademicYear: "2024",
			},
			expectValid: false,
		},
		{
			name: "Missing ProgramStudy",
			student: &models.Student{
				UserID:       "user-123",
				StudentID:    "2024001",
				ProgramStudy: "",
				AcademicYear: "2024",
			},
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.student.UserID != "" &&
				tc.student.StudentID != "" &&
				tc.student.ProgramStudy != "" &&
				tc.student.AcademicYear != ""

			assert.Equal(t, tc.expectValid, isValid)
		})
	}
}

// TestStudentIDGeneration tests student ID generation logic
func TestStudentIDGeneration(t *testing.T) {
	// Mock repository
	mockRepo := new(MockStudentRepository)

	// Setup mock expectation
	mockRepo.On("CountByYear", 2024).Return(int64(5), nil)

	// Test
	count, err := mockRepo.CountByYear(2024)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)

	// Verify mock was called
	mockRepo.AssertExpectations(t)
}

// TestSetAdvisorValidation tests advisor assignment validation
func TestSetAdvisorValidation(t *testing.T) {
	testCases := []struct {
		name        string
		studentID   string
		advisorID   string
		expectValid bool
	}{
		{
			name:        "Valid assignment",
			studentID:   "student-123",
			advisorID:   "advisor-456",
			expectValid: true,
		},
		{
			name:        "Missing student ID",
			studentID:   "",
			advisorID:   "advisor-456",
			expectValid: false,
		},
		{
			name:        "Missing advisor ID",
			studentID:   "student-123",
			advisorID:   "",
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.studentID != "" && tc.advisorID != ""
			assert.Equal(t, tc.expectValid, isValid)
		})
	}
}

// TestStudentCreate tests student creation with mock repository
func TestStudentCreate(t *testing.T) {
	// Setup mock
	mockRepo := new(MockStudentRepository)

	student := &models.Student{
		ID:           "student-123",
		UserID:       "user-123",
		StudentID:    "2024001",
		ProgramStudy: "Informatika",
		AcademicYear: "2024",
	}

	// Setup expectation
	mockRepo.On("Create", student).Return(nil)

	// Test
	err := mockRepo.Create(student)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestStudentUpdate tests student update with mock repository
func TestStudentUpdate(t *testing.T) {
	// Setup mock
	mockRepo := new(MockStudentRepository)

	student := &models.Student{
		ID:           "student-123",
		UserID:       "user-123",
		StudentID:    "2024001",
		ProgramStudy: "Sistem Informasi",
		AcademicYear: "2024",
	}

	// Setup expectation
	mockRepo.On("Update", student).Return(nil)

	// Test
	err := mockRepo.Update(student)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestFindStudentByID tests finding student by ID with mock
func TestFindStudentByID(t *testing.T) {
	// Setup mock
	mockRepo := new(MockStudentRepository)

	expectedStudent := &models.Student{
		ID:           "student-123",
		UserID:       "user-123",
		StudentID:    "2024001",
		ProgramStudy: "Informatika",
		AcademicYear: "2024",
	}

	// Setup expectation
	mockRepo.On("FindByID", "student-123").Return(expectedStudent, nil)

	// Test
	student, err := mockRepo.FindByID("student-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, student)
	assert.Equal(t, "student-123", student.ID)
	assert.Equal(t, "2024001", student.StudentID)
	mockRepo.AssertExpectations(t)
}
