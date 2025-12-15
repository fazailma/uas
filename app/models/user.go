package models

import "time"

type User struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	RoleID       string    `json:"role_id"`
	IsActive     bool      `json:"is_active"`
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

// CreateUserRequest represents the request payload for creating a new user
type CreateUserRequest struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	FullName     string `json:"full_name" validate:"required"`
	RoleID       string `json:"role_id" validate:"required"`
	StudentID    string `json:"student_id,omitempty"`
	ProgramStudy string `json:"program_study,omitempty"`
	AcademicYear string `json:"academic_year,omitempty"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
	RoleID   string `json:"role_id,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// UserResponse represents the response format for user data
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	RoleID   string `json:"role_id"`
	IsActive bool   `json:"is_active"`
}
