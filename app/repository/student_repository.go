package repository

import (
	"fmt"

	"UAS/app/models"
	"UAS/database"
)

// StudentRepository handles student database operations
type StudentRepository struct{}

// NewStudentRepository creates a new instance of StudentRepository
func NewStudentRepository() *StudentRepository {
	return &StudentRepository{}
}

// Create creates a new student record
func (r *StudentRepository) Create(student *models.Student) error {
	return database.DB.Create(student).Error
}

// FindByUserID finds student by user ID
func (r *StudentRepository) FindByUserID(userID string) (*models.Student, error) {
	var student models.Student
	err := database.DB.Where("user_id = ?", userID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// FindByStudentID finds student by student ID (NIM)
func (r *StudentRepository) FindByStudentID(studentID string) (*models.Student, error) {
	var student models.Student
	err := database.DB.Where("student_id = ?", studentID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// Update updates a student record
func (r *StudentRepository) Update(id string, student *models.Student) error {
	return database.DB.Model(&models.Student{}).Where("id = ?", id).Updates(student).Error
}

// Delete deletes a student record
func (r *StudentRepository) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.Student{}).Error
}

// FindByAdvisorID finds all students guided by an advisor
func (r *StudentRepository) FindByAdvisorID(advisorID string) ([]models.Student, error) {
	var students []models.Student
	err := database.DB.Where("advisor_id = ?", advisorID).Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

// CountByYear counts students with student_id starting with year
func (r *StudentRepository) CountByYear(year int) (int64, error) {
	var count int64
	yearStr := fmt.Sprintf("%d%%", year)
	err := database.DB.Model(&models.Student{}).
		Where("student_id LIKE ?", yearStr).
		Count(&count).Error
	return count, err
}

// CountByAdvisorID counts students assigned to an advisor
func (r *StudentRepository) CountByAdvisorID(advisorID string) (int64, error) {
	var count int64
	err := database.DB.Model(&models.Student{}).
		Where("advisor_id = ?", advisorID).
		Count(&count).Error
	return count, err
}
