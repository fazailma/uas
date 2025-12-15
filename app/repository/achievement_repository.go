package repository

import (
	"time"

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
	err := database.DB.Where("id = ? AND deleted_at IS NULL", id).First(&achievement).Error
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

// FindByStudentID finds all achievements by student ID
func (r *AchievementRepository) FindByStudentID(studentID string) ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Where("student_id = ? AND deleted_at IS NULL", studentID).Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// FindDraftByStudentID finds all draft achievements by student ID
func (r *AchievementRepository) FindDraftByStudentID(studentID string) ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Where("student_id = ? AND status = ? AND deleted_at IS NULL", studentID, "draft").
		Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// FindAll finds all achievements (admin view)
func (r *AchievementRepository) FindAll() ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Where("deleted_at IS NULL").Order("created_at DESC").Find(&achievements).Error
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
	return database.DB.Model(&models.AchievementReference{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("status", status).Error
}

// Delete soft delete an achievement (FR-005)
func (r *AchievementRepository) Delete(id string) error {
	now := time.Now()
	return database.DB.Model(&models.AchievementReference{}).
		Where("id = ?", id).
		Update("deleted_at", now).Error
}

// FindByStatus finds achievements by status
func (r *AchievementRepository) FindByStatus(status string) ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Where("status = ? AND deleted_at IS NULL", status).Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, err
	}
	return achievements, nil
}

// FindByStudentIDs finds achievements by multiple student IDs
func (r *AchievementRepository) FindByStudentIDs(studentIDs []string) ([]models.AchievementReference, error) {
	var achievements []models.AchievementReference
	err := database.DB.Where("student_id IN ? AND deleted_at IS NULL", studentIDs).Order("created_at DESC").Find(&achievements).Error
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

// CountByStatus returns count of achievements grouped by status
func (r *AchievementRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := database.DB.Model(&models.AchievementReference{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}

	return counts, nil
}

// CountTotal returns total count of achievements
func (r *AchievementRepository) CountTotal() (int64, error) {
	var total int64
	err := database.DB.Model(&models.AchievementReference{}).Count(&total).Error
	return total, err
}

// FindAllWithPagination finds all achievements with pagination
func (r *AchievementRepository) FindAllWithPagination(page, pageSize int) ([]models.AchievementReference, int64, error) {
	var achievements []models.AchievementReference
	var total int64

	offset := (page - 1) * pageSize

	if err := database.DB.Model(&models.AchievementReference{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := database.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&achievements).Error
	if err != nil {
		return nil, 0, err
	}

	return achievements, total, nil
}
