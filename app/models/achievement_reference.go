package models

import "time"

type AchievementReference struct {
	ID                 string    `json:"id" gorm:"primaryKey"`
	StudentID          string    `json:"student_id"`
	MongoAchievementID string    `json:"mongo_achievement_id"`
	Status             string    `json:"status"` // draft, submitted, verified, rejected
	SubmittedAt        time.Time `json:"submitted_at"`
	VerifiedAt         time.Time `json:"verified_at"`
	VerifiedBy         string    `json:"verified_by"`
	RejectionNote      string    `json:"rejection_note"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// AchievementCreateRequest represents request to create achievement (FR-003)
type AchievementCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"required"`
	Date        string `json:"date" binding:"required"` // Format: YYYY-MM-DD
	ProofURL    string `json:"proof_url"`
}

// AchievementSubmitRequest represents request to submit achievement (FR-004)
type AchievementSubmitRequest struct {
	AchievementID string `json:"achievement_id" binding:"required"`
}

// AchievementDeleteRequest represents request to delete achievement (FR-005)
type AchievementDeleteRequest struct {
	AchievementID string `json:"achievement_id" binding:"required"`
}

// AchievementUpdateRequest represents request to update achievement
type AchievementUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Date        string `json:"date"` // Format: YYYY-MM-DD
	ProofURL    string `json:"proof_url"`
}
