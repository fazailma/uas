package database

import (
	"log"

	"UAS/app/models"

	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")

	// AutoMigrate akan membuat tabel jika belum ada
	// 7 model: user, role, permission, role_permission, student, lecturer, achievement_reference
	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.Student{},
		&models.Lecturer{},
		&models.AchievementReference{},
	)

	if err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	log.Println("Migrations completed successfully")
}
