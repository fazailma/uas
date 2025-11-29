package models

import "time"

type Lecturer struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	UserID     string    `json:"user_id"`
	LecturerID string    `json:"lecturer_id"` // NIP
	Department string    `json:"department"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
