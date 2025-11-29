package models

type Permission struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}
