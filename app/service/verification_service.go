package service

import (
	"context"
	"errors"
	"time"

	"UAS/app/repository"
)

type VerificationService struct {
	achievementRepo *repository.AchievementRepository
	studentRepo     *repository.StudentRepository
	mongoRepo       *repository.MongoAchievementRepository
}

func NewVerificationService() *VerificationService {
	return &VerificationService{
		achievementRepo: repository.NewAchievementRepository(),
		studentRepo:     repository.NewStudentRepository(),
		mongoRepo:       repository.NewMongoAchievementRepository(),
	}
}

// GetGuidedStudentsAchievements gets all achievements from students guided by this lecturer
// FR-006: View Prestasi Mahasiswa Bimbingan
func (s *VerificationService) GetGuidedStudentsAchievements(dosenID string) ([]map[string]interface{}, error) {
	// Step 1: Find lecturer by user ID
	lecturer, err := repository.NewLecturerRepository().FindByUserID(dosenID)
	if err != nil {
		return nil, errors.New("lecturer not found")
	}

	// Step 2: Get all students guided by this lecturer
	students, err := s.studentRepo.FindByAdvisorID(lecturer.ID)
	if err != nil {
		return nil, errors.New("failed to fetch guided students: " + err.Error())
	}

	if len(students) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Step 3: Get student IDs
	var studentIDs []string
	for _, student := range students {
		studentIDs = append(studentIDs, student.UserID)
	}

	// Step 4: Get achievement references for these students
	achievements, err := s.achievementRepo.FindByStudentIDs(studentIDs)
	if err != nil {
		return nil, errors.New("failed to fetch achievements: " + err.Error())
	}

	// Step 5: Fetch MongoDB details and combine
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result []map[string]interface{}
	for _, ach := range achievements {
		// Fetch MongoDB document
		mongoAch, err := s.mongoRepo.FindByID(ctx, ach.MongoAchievementID)
		if err != nil {
			continue // Skip if MongoDB document not found
		}

		result = append(result, map[string]interface{}{
			"id":             ach.ID,
			"student_id":     ach.StudentID,
			"title":          mongoAch.Title,
			"description":    mongoAch.Description,
			"category":       mongoAch.Category,
			"date":           mongoAch.Date,
			"proof_url":      mongoAch.ProofURL,
			"status":         ach.Status,
			"submitted_at":   ach.SubmittedAt,
			"verified_at":    ach.VerifiedAt,
			"verified_by":    ach.VerifiedBy,
			"rejection_note": ach.RejectionNote,
			"created_at":     ach.CreatedAt,
		})
	}

	return result, nil
}

// VerifyAchievement verifies an achievement by dosen wali
// FR-007: Verify Prestasi
func (s *VerificationService) VerifyAchievement(achievementID, dosenID string) error {
	// Get achievement
	achievement, err := s.achievementRepo.FindByID(achievementID)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Verify ownership - check if dosen is advisor of this student
	student, err := s.studentRepo.FindByUserID(achievement.StudentID)
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
	err = s.achievementRepo.VerifyAchievement(achievementID, dosenID)
	if err != nil {
		return errors.New("failed to verify achievement: " + err.Error())
	}

	// TODO: Create notification for student

	return nil
}

// RejectAchievement rejects an achievement with rejection note
// FR-008: Reject Prestasi
func (s *VerificationService) RejectAchievement(achievementID, rejectionNote string) error {
	// Get achievement
	achievement, err := s.achievementRepo.FindByID(achievementID)
	if err != nil {
		return errors.New("achievement not found")
	}

	// Check status
	if achievement.Status != "submitted" {
		return errors.New("only submitted achievements can be rejected")
	}

	// Update status to rejected with note
	err = s.achievementRepo.RejectAchievement(achievementID, rejectionNote)
	if err != nil {
		return errors.New("failed to reject achievement: " + err.Error())
	}

	// TODO: Create notification for student with rejection reason

	return nil
}

// GetGuidedStudentsAchievementsWithRoleCheck gets achievements with role validation
func (s *VerificationService) GetGuidedStudentsAchievementsWithRoleCheck(dosenID, role string) ([]map[string]interface{}, error) {
	if role != "Dosen Wali" {
		return nil, errors.New("only dosen wali can access this endpoint")
	}
	return s.GetGuidedStudentsAchievements(dosenID)
}

// VerifyAchievementWithRoleCheck verifies achievement with role validation
func (s *VerificationService) VerifyAchievementWithRoleCheck(achievementID, dosenID, role string) error {
	if role != "Dosen Wali" {
		return errors.New("only dosen wali can verify achievements")
	}
	return s.VerifyAchievement(achievementID, dosenID)
}

// RejectAchievementWithRoleCheck rejects achievement with role validation
func (s *VerificationService) RejectAchievementWithRoleCheck(achievementID, rejectionNote, role string) error {
	if role != "Dosen Wali" {
		return errors.New("only dosen wali can reject achievements")
	}
	return s.RejectAchievement(achievementID, rejectionNote)
}
