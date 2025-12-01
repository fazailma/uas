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

// Create creates a new achievement (FR-003)
func (r *AchievementRepository) Create(achievement *models.AchievementReference) error {
	return database.DB.Create(achievement).Error
}

// FindByID finds achievement by ID
func (r *AchievementRepository) FindByID(id string) (*models.AchievementReference, error) {
	var achievement models.AchievementReference
	err := database.DB.Where("id = ?", id).First(&achievement).Error
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// FindByStudentID finds all achievements by student ID
func (r *AchievementRepository) FindByStudentID(studentID string) ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Where("student_id = ?", studentID).Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// FindDraftByStudentID finds all draft achievements by student ID
func (r *AchievementRepository) FindDraftByStudentID(studentID string) ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Where("student_id = ? AND status = ?", studentID, "draft").
		Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// FindAll finds all achievements (admin view)
func (r *AchievementRepository) FindAll() ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// Update updates an achievement
func (r *AchievementRepository) Update(id string, achievement *models.AchievementReference) error {
	return database.DB.Model(&models.AchievementReference{}).Where("id = ?", id).Updates(achievement).Error
}

// UpdateStatus updates achievement status (FR-004)
func (r *AchievementRepository) UpdateStatus(id string, status string) error {
	return database.DB.Model(&models.AchievementReference{}).Where("id = ?", id).Update("status", status).Error
}

// Delete soft delete an achievement (FR-005)
func (r *AchievementRepository) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.AchievementReference{}).Error
}

// FindByStatus finds achievements by status
func (r *AchievementRepository) FindByStatus(status string) ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Where("status = ?", status).Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// FindByStudentIDs finds achievements by multiple student IDs
func (r *AchievementRepository) FindByStudentIDs(studentIDs []string) ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Where("student_id IN ?", studentIDs).Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// VerifyAchievement verifies an achievement (FR-007)
func (r *AchievementRepository) VerifyAchievement(id string, verifiedBy string) error {
	return database.DB.Model(&models.AchievementReference{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "verified",
			"verified_by": verifiedBy,
			"verified_at": database.DB.Statement.DB.NowFunc(),
		}).Error
}

// RejectAchievement rejects an achievement with note (FR-008)
func (r *AchievementRepository) RejectAchievement(id string, rejectionNote string) error {
	return database.DB.Model(&models.AchievementReference{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         "rejected",
			"rejection_note": rejectionNote,
		}).Error
}
