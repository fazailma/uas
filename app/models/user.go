package models

import "time"

type User struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex"`
	Email        string    `json:"email" gorm:"uniqueIndex"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	RoleID       string    `json:"role_id"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// LoginCredential represents login credential request
type LoginCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token   string `json:"token"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}
