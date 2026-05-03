// app/repository/permission.go

package repository

import (
	"context"
	"errors"
	"iam-box/app/entities"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

// Create - grant a new permission
func (r *PermissionRepository) Create(ctx context.Context, p *entities.Permission) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// Repository layer — returns the permission, not just bool
func (r *PermissionRepository) Find(ctx context.Context, userID string, action entities.Action, resourceType string, resourceID *string) (*entities.Permission, error) {
	var perm entities.Permission

	query := r.db.WithContext(ctx).
		Where("user_id = ? AND action = ? AND resource_type = ?", userID, action, resourceType)

	if resourceID != nil {
		query = query.Where("resource_id = ?", *resourceID)
	} else {
		query = query.Where("resource_id IS NULL")
	}

	err := query.First(&perm).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Not found, no error
	}

	return &perm, err
}

// List - paginated permissions
func (r *PermissionRepository) List(ctx context.Context, limit, offset int) ([]entities.Permission, error) {
	var perms []entities.Permission
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("granted_at DESC").Find(&perms).Error
	return perms, err
}

// GetByID - find permission by ID
func (r *PermissionRepository) GetByID(ctx context.Context, id uint) (*entities.Permission, error) {
	var perm entities.Permission
	err := r.db.WithContext(ctx).First(&perm, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &perm, err
}

// GetByUser - find all permissions for a user
func (r *PermissionRepository) GetByUser(ctx context.Context, userID string) ([]entities.Permission, error) {
	var perms []entities.Permission
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&perms).Error
	return perms, err
}

// Check - the key processing unit
// it allows or denies the access to the resource
// and also logs it as a decision
// this one is also potentially heavy loaded
// and should take advantage of the caching (Redis)
func (r *PermissionRepository) Check(ctx context.Context, userID, resourceType string, action entities.Action, resourceID *string) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&entities.Permission{}).
		Where("user_id = ? AND action = ? AND resource_type = ?", userID, action, resourceType)

	if resourceID != nil {
		// Check for specific OR wildcard
		query = query.Where("resource_id = ? OR resource_id IS NULL", *resourceID)
	} else {
		// Check for wildcard only
		query = query.Where("resource_id IS NULL")
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// Delete - remove a permission
func (r *PermissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Permission{}, id).Error
}
