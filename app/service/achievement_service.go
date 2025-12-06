package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"UAS/app/models"
	"UAS/app/repository"
)

type AchievementService struct {
	pgRepo    *repository.AchievementRepository
	mongoRepo *repository.MongoAchievementRepository
}

func NewAchievementService() *AchievementService {
	return &AchievementService{
		pgRepo:    repository.NewAchievementRepository(),
		mongoRepo: repository.NewMongoAchievementRepository(),
	}
}

// CreateAchievement creates achievement in both MongoDB and PostgreSQL
// FR-003: Submit Prestasi
func (s *AchievementService) CreateAchievement(userID, role string, req models.AchievementCreateRequest) (*models.AchievementReference, error) {
	// Validate required fields
	if req.Title == "" || req.Category == "" || req.Date == "" {
		return nil, errors.New("title, category, and date are required")
	}

	// Only Mahasiswa can create
	if role != "Mahasiswa" {
		return nil, errors.New("only mahasiswa can submit achievements")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Create MongoDB document with full achievement details
	mongoAchievement := &models.MongoAchievement{
		StudentID:   userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Date:        req.Date,
		ProofURL:    req.ProofURL,
	}

	mongoAch, err := s.mongoRepo.Create(ctx, mongoAchievement)
	if err != nil {
		return nil, errors.New("failed to save achievement to MongoDB: " + err.Error())
	}

	// 2. Create PostgreSQL reference
	pgAchievement := &models.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          userID,
		MongoAchievementID: mongoAch.ID.Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	err = s.pgRepo.Create(pgAchievement)
	if err != nil {
		// Rollback MongoDB if PostgreSQL fails
		s.mongoRepo.SoftDelete(ctx, mongoAch.ID.Hex())
		return nil, errors.New("failed to save achievement reference to PostgreSQL: " + err.Error())
	}

	return pgAchievement, nil
}

// SubmitAchievement updates achievement status to 'submitted'
// FR-004: Submit untuk Verifikasi
func (s *AchievementService) SubmitAchievement(id, userID, role string) error {
	achievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Ownership check - only mahasiswa can submit their own
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return errors.New("you can only submit your own achievements")
	}

	// Can only submit draft
	if achievement.Status != "draft" {
		return errors.New("only draft achievements can be submitted")
	}

	// Update status to submitted in PostgreSQL
	err = s.pgRepo.UpdateStatus(id, "submitted")
	if err != nil {
		return err
	}

	// TODO: Create notification for dosen wali

	return nil
}

// UpdateAchievement updates achievement details (only draft)
func (s *AchievementService) UpdateAchievement(id, userID, role string, req models.AchievementUpdateRequest) (*models.AchievementReference, error) {
	pgAchievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("achievement not found")
	}

	// Ownership check - only mahasiswa can update their own
	if role == "Mahasiswa" && pgAchievement.StudentID != userID {
		return nil, errors.New("you can only update your own achievements")
	}

	// Can only update draft achievements
	if pgAchievement.Status != "draft" {
		return nil, errors.New("only draft achievements can be updated")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update MongoDB document
	mongoID := pgAchievement.MongoAchievementID
	mongoAch := &models.MongoAchievement{
		StudentID:   userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Date:        req.Date,
		ProofURL:    req.ProofURL,
	}

	_, err = s.mongoRepo.Update(ctx, mongoID, mongoAch)
	if err != nil {
		return nil, errors.New("failed to update achievement in MongoDB: " + err.Error())
	}

	// Update PostgreSQL reference timestamp
	pgAchievement.UpdatedAt = time.Now()
	err = s.pgRepo.Update(id, pgAchievement)
	if err != nil {
		return nil, errors.New("failed to update achievement reference in PostgreSQL: " + err.Error())
	}

	return pgAchievement, nil
}

// DeleteAchievement soft deletes achievement from both MongoDB and marks PostgreSQL
// FR-005: Hapus Prestasi
func (s *AchievementService) DeleteAchievement(id, userID, role string) error {
	pgAchievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Ownership check - only mahasiswa can delete their own
	if role == "Mahasiswa" && pgAchievement.StudentID != userID {
		return errors.New("you can only delete your own achievements")
	}

	// Can only delete draft
	if pgAchievement.Status != "draft" {
		return errors.New("only draft achievements can be deleted")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Soft delete MongoDB document
	mongoID := pgAchievement.MongoAchievementID
	err = s.mongoRepo.SoftDelete(ctx, mongoID)
	if err != nil {
		return errors.New("failed to delete achievement from MongoDB: " + err.Error())
	}

	// 2. Soft delete PostgreSQL reference
	err = s.pgRepo.Delete(id)
	if err != nil {
		return errors.New("failed to delete achievement reference from PostgreSQL: " + err.Error())
	}

	return nil
}

// VerifyAchievement verifies an achievement by dosen wali
// FR-007: Verify Prestasi
func (s *AchievementService) VerifyAchievement(achievementID, dosenID string) error {
	// Get achievement
	achievement, err := s.pgRepo.FindByID(achievementID)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Verify ownership - check if dosen is advisor of this student
	student, err := repository.NewStudentRepository().FindByUserID(achievement.StudentID)
	if err != nil {
		return errors.New("student not found")
	}

	lecturer, err := repository.NewLecturerRepository().FindByUserID(dosenID)
	if err != nil {
		return errors.New("lecturer not found")
	}

	if student.AdvisorID != lecturer.ID {
		return errors.New("you can only verify achievements of your guided students")
	}

	// Check status
	if achievement.Status != "submitted" {
		return errors.New("only submitted achievements can be verified")
	}

	// Update status to verified
	err = s.pgRepo.VerifyAchievement(achievementID, dosenID)
	if err != nil {
		return errors.New("failed to verify achievement: " + err.Error())
	}

	// TODO: Create notification for student

	return nil
}

// RejectAchievement rejects an achievement with rejection note
// FR-008: Reject Prestasi
func (s *AchievementService) RejectAchievement(achievementID, rejectionNote string) error {
	// Get achievement
	achievement, err := s.pgRepo.FindByID(achievementID)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Check status
	if achievement.Status != "submitted" {
		return errors.New("only submitted achievements can be rejected")
	}

	// Update status to rejected with note
	err = s.pgRepo.RejectAchievement(achievementID, rejectionNote)
	if err != nil {
		return errors.New("failed to reject achievement: " + err.Error())
	}

	// TODO: Create notification for student with rejection reason

	return nil
}

// ListAchievements lists achievements based on user role
func (s *AchievementService) ListAchievements(userID, role string) ([]models.AchievementReference, error) {
	switch role {
	case "Admin":
		return s.pgRepo.FindAll()
	case "Mahasiswa":
		return s.pgRepo.FindByStudentID(userID)
	case "Dosen Wali":
		// TODO: Implement logic to find students under guidance
		return []models.AchievementReference{}, nil
	default:
		return s.pgRepo.FindAll()
	}
}

// GetAchievementDetail gets achievement details with ownership check
func (s *AchievementService) GetAchievementDetail(id, userID, role string) (*models.AchievementReference, error) {
	achievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Ownership check for Mahasiswa
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return nil, errors.New("you can only view your own achievements")
	}

	return achievement, nil
}

// GetAchievementHistory retrieves achievement status history - pure logic
func (s *AchievementService) GetAchievementHistory(id, userID, role string) (map[string]interface{}, error) {
	achievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("achievement not found")
	}

	// Ownership check for Mahasiswa
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return nil, errors.New("you can only view your own achievement history")
	}

	// Build history timeline
	history := map[string]interface{}{
		"id":     id,
		"status": achievement.Status,
		"timeline": []map[string]interface{}{
			{
				"status":    "draft",
				"timestamp": achievement.CreatedAt,
				"message":   "Achievement created",
			},
		},
	}

	// Add submitted event if applicable
	if !achievement.SubmittedAt.IsZero() {
		timeline := history["timeline"].([]map[string]interface{})
		timeline = append(timeline, map[string]interface{}{
			"status":    "submitted",
			"timestamp": achievement.SubmittedAt,
			"message":   "Achievement submitted for verification",
		})
		history["timeline"] = timeline
	}

	// Add verified event if applicable
	if !achievement.VerifiedAt.IsZero() {
		timeline := history["timeline"].([]map[string]interface{})
		timeline = append(timeline, map[string]interface{}{
			"status":      "verified",
			"timestamp":   achievement.VerifiedAt,
			"verified_by": achievement.VerifiedBy,
			"message":     "Achievement verified",
		})
		history["timeline"] = timeline
	}

	// Add rejected event if applicable
	if achievement.Status == "rejected" && achievement.RejectionNote != "" {
		timeline := history["timeline"].([]map[string]interface{})
		timeline = append(timeline, map[string]interface{}{
			"status":  "rejected",
			"message": achievement.RejectionNote,
		})
		history["timeline"] = timeline
	}

	return history, nil
}

// ValidateAchievementOwnership checks if user owns the achievement - pure logic
func (s *AchievementService) ValidateAchievementOwnership(id, userID string) error {
	achievement, err := s.pgRepo.FindByID(id)
	if err != nil {
		return errors.New("achievement not found")
	}

	if achievement.StudentID != userID {
		return errors.New("you can only access your own achievements")
	}

	return nil
}
