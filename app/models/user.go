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

// RegisterRequest represents register request with additional fields
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	RoleID   string `json:"role_id"`
}

// UserProfile represents user profile in response
type UserProfile struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"fullName"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// LoginResponseData represents data inside login response
type LoginResponseData struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refreshToken"`
	User         UserProfile `json:"user"`
}

// LoginResponse represents login response wrapper
type LoginResponse struct {
	Status string            `json:"status"`
	Data   LoginResponseData `json:"data"`
}

// RegisterResponse represents register response
type RegisterResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}
