package repository

import (
	"UAS/app/models"
	"UAS/database"
)

// LecturerRepository handles lecturer database operations
type LecturerRepository struct{}

// NewLecturerRepository creates a new instance of LecturerRepository
func NewLecturerRepository() *LecturerRepository {
	return &LecturerRepository{}
}

// Create creates a new lecturer record
func (r *LecturerRepository) Create(lecturer *models.Lecturer) error {
	return database.DB.Create(lecturer).Error
}

// FindByUserID finds lecturer by user ID
func (r *LecturerRepository) FindByUserID(userID string) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	err := database.DB.Where("user_id = ?", userID).First(&lecturer).Error
	if err != nil {
		return nil, err
	}
	return &lecturer, nil
}

// FindByLecturerID finds lecturer by lecturer ID (NIP)
func (r *LecturerRepository) FindByLecturerID(lecturerID string) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	err := database.DB.Where("lecturer_id = ?", lecturerID).First(&lecturer).Error
	if err != nil {
		return nil, err
	}
	return &lecturer, nil
}

// Update updates a lecturer record
func (r *LecturerRepository) Update(id string, lecturer *models.Lecturer) error {
	return database.DB.Model(&models.Lecturer{}).Where("id = ?", id).Updates(lecturer).Error
}

// Delete deletes a lecturer record
func (r *LecturerRepository) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.Lecturer{}).Error
}

// FindAll finds all lecturers
func (r *LecturerRepository) FindAll() ([]models.Lecturer, error) {
	var lecturers []models.Lecturer
	err := database.DB.Find(&lecturers).Error
	return lecturers, err
}

// CountTotal counts total lecturers
func (r *LecturerRepository) CountTotal() (int64, error) {
	var count int64
	err := database.DB.Model(&models.Lecturer{}).Count(&count).Error
	return count, err
}

// FindByID finds lecturer by ID
func (r *LecturerRepository) FindByID(id string) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	err := database.DB.Where("id = ?", id).First(&lecturer).Error
	if err != nil {
		return nil, err
	}
	return &lecturer, nil
}
