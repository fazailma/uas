package repository

import (
	"UAS/app/models"
	"UAS/database"
)

// RoleRepository handles role database operations
type RoleRepository struct{}

// NewRoleRepository creates a new instance of RoleRepository
func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

// FindByID finds role by id
func (r *RoleRepository) FindByID(id string) (*models.Role, error) {
	var role models.Role
	err := database.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName finds role by name
func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := database.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Create creates a new role
func (r *RoleRepository) Create(role *models.Role) error {
	return database.DB.Create(role).Error
}

// Update updates a role
func (r *RoleRepository) Update(role *models.Role) error {
	return database.DB.Save(role).Error
}

// FindAll retrieves all roles
func (r *RoleRepository) FindAll() ([]*models.Role, error) {
	var roles []*models.Role
	err := database.DB.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// Delete deletes a role
func (r *RoleRepository) Delete(id string) error {
	return database.DB.Where("id = ?", id).Delete(&models.Role{}).Error
}

// PermissionRepository handles permission database operations
type PermissionRepository struct{}

// NewPermissionRepository creates a new instance of PermissionRepository
func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{}
}

// FindByID finds permission by id
func (r *PermissionRepository) FindByID(id string) (*models.Permission, error) {
	var permission models.Permission
	err := database.DB.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// FindByName finds permission by name
func (r *PermissionRepository) FindByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := database.DB.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// Create creates a new permission
func (r *PermissionRepository) Create(permission *models.Permission) error {
	return database.DB.Create(permission).Error
}

// RolePermissionRepository handles role permission database operations
type RolePermissionRepository struct{}

// NewRolePermissionRepository creates a new instance of RolePermissionRepository
func NewRolePermissionRepository() *RolePermissionRepository {
	return &RolePermissionRepository{}
}

// AssignPermissionToRole assigns a permission to a role
func (r *RolePermissionRepository) AssignPermissionToRole(roleID string, permissionID string) error {
	rolePermission := models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return database.DB.Create(&rolePermission).Error
}

// GetPermissionsByRole gets all permissions for a role
func (r *RolePermissionRepository) GetPermissionsByRole(roleID string) ([]models.Permission, error) {
	// If roleID is empty, return empty list (avoid invalid uuid queries)
	if roleID == "" {
		return []models.Permission{}, nil
	}

	var permissions []models.Permission
	err := database.DB.
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
