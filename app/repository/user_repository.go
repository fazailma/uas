package repository

import (
	"UAS/app/models"
	"UAS/database"
)

// UserRepository handles user database operations
type UserRepository struct{}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByUsername finds user by username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID finds user by id
func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserWithRoleAndPermissions gets user with role and permissions
func (r *UserRepository) GetUserWithRoleAndPermissions(userID string) (*models.User, []models.Permission, error) {
	var user models.User
	err := database.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, nil, err
	}

	// Get permissions for this role
	// If user has no role assigned, return empty permissions (avoid invalid uuid queries)
	if user.RoleID == "" {
		var empty []models.Permission
		return &user, empty, nil
	}

	var permissions []models.Permission
	err = database.DB.
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", user.RoleID).
		Find(&permissions).Error
	if err != nil {
		return nil, nil, err
	}

	return &user, permissions, nil
}

// Create creates a new user record
func (r *UserRepository) Create(user *models.User) error {
	// If RoleID is empty, omit it so DB can use NULL/default (avoid inserting empty string into uuid column)
	if user.RoleID == "" {
		return database.DB.Omit("role_id").Create(user).Error
	}
	return database.DB.Create(user).Error
}

// Update updates a user record
func (r *UserRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}

// FindAll retrieves all users with pagination
func (r *UserRepository) FindAll(page, pageSize int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := database.DB.Where("is_active = ?", true)

	if err := query.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
