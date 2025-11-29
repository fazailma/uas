package repository

import (
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
