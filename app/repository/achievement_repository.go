package repository

import (
	"UAS/app/models"
	"UAS/database"
)

// AchievementRepository handles achievement database operations
type AchievementRepository struct{}

// NewAchievementRepository creates a new instance of AchievementRepository
func NewAchievementRepository() *AchievementRepository {
	return &AchievementRepository{}
}

// Create creates a new achievement
func (r *AchievementRepository) Create(achievement *models.Achievement) error {
	return database.DB.Create(achievement).Error
}

// FindByID finds achievement by ID
func (r *AchievementRepository) FindByID(id string) (*models.Achievement, error) {
	var achievement models.Achievement
	err := database.DB.Where("id = ?", id).First(&achievement).Error
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// FindByStudentID finds all achievements by student ID
func (r *AchievementRepository) FindByStudentID(studentID string) ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := database.DB.Where("student_id = ?", studentID).Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// FindAll finds all achievements (admin view)
func (r *AchievementRepository) FindAll() ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := database.DB.Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// Update updates an achievement
func (r *AchievementRepository) Update(id string, achievement *models.Achievement) error {
	return database.DB.Model(&models.Achievement{}).Where("id = ?", id).Updates(achievement).Error
}

// UpdateStatus updates achievement status
func (r *AchievementRepository) UpdateStatus(id string, status string) error {
	return database.DB.Model(&models.Achievement{}).Where("id = ?", id).Update("status", status).Error
}

// Delete deletes an achievement
func (r *AchievementRepository) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.Achievement{}).Error
}

// FindByStatus finds achievements by status
func (r *AchievementRepository) FindByStatus(status string) ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := database.DB.Where("status = ?", status).Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}
