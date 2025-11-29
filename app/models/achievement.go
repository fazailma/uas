package models

import "time"

type Achievement struct {
	ID            string    `json:"id"`
	StudentID     string    `json:"student_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Date          time.Time `json:"date"`
	ProofURL      string    `json:"proof_url"`
	Status        string    `json:"status"`
	Points        int       `json:"points"`
	RejectionNote string    `json:"rejection_note"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AchievementCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Date        string `json:"date"`
	ProofURL    string `json:"proof_url"`
}

type AchievementUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Date        string `json:"date"`
	ProofURL    string `json:"proof_url"`
}

type AchievementVerifyRequest struct {
	Points int `json:"points"`
}

type AchievementRejectRequest struct {
	Reason string `json:"reason"`
}
