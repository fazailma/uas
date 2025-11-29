package models

type RolePermission struct {
	RoleID       string `gorm:"primaryKey"`
	PermissionID string `gorm:"primaryKey"`
}
