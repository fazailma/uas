package models

// CreateUserRequest represents the request payload for creating a new user
type CreateUserRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	RoleID       string `json:"role_id"`
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

// AchievementDetailResponse represents the response format for achievement data
type AchievementDetailResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Date        string `json:"date"`
	ProofURL    string `json:"proof_url"`
	Status      string `json:"status"`
	StudentID   string `json:"student_id"`
	CreatedAt   string `json:"created_at"`
}

// CreateAchievementRequest represents request to create achievement
type CreateAchievementRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"required"`
	Date        string `json:"date" validate:"required"`
	ProofURL    string `json:"proof_url"`
}

// UpdateAchievementRequest represents request to update achievement
type UpdateAchievementRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Date        string `json:"date"`
	ProofURL    string `json:"proof_url"`
}
