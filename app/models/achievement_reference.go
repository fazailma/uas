package models

import "time"

type AchievementReference struct {
	ID                 string     `json:"id" gorm:"primaryKey"`
	StudentID          string     `json:"student_id"`
	MongoAchievementID string     `json:"mongo_achievement_id"`
	Status             string     `json:"status"` // draft, submitted, verified, rejected
	SubmittedAt        time.Time  `json:"submitted_at"`
	VerifiedAt         time.Time  `json:"verified_at"`
	VerifiedBy         string     `json:"verified_by"`
	RejectionNote      string     `json:"rejection_note"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
}

// CreateAchievementRequest represents request to create achievement
type CreateAchievementRequest struct {
	Title           string                 `json:"title" validate:"required"`
	Description     string                 `json:"description"`
	AchievementType string                 `json:"achievement_type" validate:"required"` // 'academic', 'competition', 'organization', 'publication', 'certification', 'other'
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
}

// UpdateAchievementRequest represents request to update achievement
type UpdateAchievementRequest struct {
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	AchievementType string                 `json:"achievement_type"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
}

// AchievementDetailResponse represents the response format for achievement data
type AchievementDetailResponse struct {
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	AchievementType string                 `json:"achievement_type"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
	Status          string                 `json:"status"`
	StudentID       string                 `json:"student_id"`
	CreatedAt       string                 `json:"created_at"`
}

// SubmitAchievementRequest represents request to submit achievement
type SubmitAchievementRequest struct {
	AchievementID string `json:"achievement_id" validate:"required"`
}

// RejectAchievementRequest represents request to reject achievement
type RejectAchievementRequest struct {
	RejectionNote string `json:"rejection_note" validate:"required"`
}
