package models

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// LoginResponseWrapper wraps login response
type LoginResponseWrapper struct {
	Status string            `json:"status"`
	Data   LoginResponseData `json:"data"`
}

// TokenResponse represents token response
type TokenResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         UserProfile `json:"user"`
}

// StatisticsResponse represents achievement statistics
type StatisticsResponse struct {
	TotalAchievements int            `json:"total_achievements"`
	ByStatus          map[string]int `json:"by_status"`
	ByType            map[string]int `json:"by_type"`
	TotalPoints       int            `json:"total_points"`
}

// AchievementListResponse represents a list of achievements
type AchievementListResponse struct {
	Achievements []AchievementReference `json:"achievements"`
	Total        int                    `json:"total"`
}

// StudentListResponse represents a list of students
type StudentListResponse struct {
	Students []Student `json:"students"`
	Total    int       `json:"total"`
}

// LecturerListResponse represents a list of lecturers
type LecturerListResponse struct {
	Lecturers []Lecturer `json:"lecturers"`
	Total     int        `json:"total"`
}

// UserListResponse represents a list of users
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}
