package models

// ErrorResponse represents an error with HTTP status and message
type ErrorResponse struct {
	Status  int
	Message string
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	return e.Message
}
