// app/service/permission.go

package service

import (
	"context"
	"fmt"
	"iam-box/app/entities"
	"iam-box/app/repository"
	"time"
)

type PermissionSvc struct {
	ID           uint
	UserID       string
	Action       entities.Action
	ResourceType string
	ResourceID   *string
	GrantedAt    time.Time
	CreatedBy    *string
}

type PermissionService struct {
	permissionRepository repository.PermissionRepository
	decisionRepository   repository.DecisionRepository
}

func NewPermissionService(
	permissionRepository repository.PermissionRepository,
	decisionRepository repository.DecisionRepository,
) *PermissionService {
	return &PermissionService{
		permissionRepository: permissionRepository,
		decisionRepository:   decisionRepository,
	}
}

func (s *PermissionService) Create(
	ctx context.Context,
	userID, action, resourceType string,
	resourceID *string,
) error {
	a := entities.Action(action)
	if !a.Valid() {
		return fmt.Errorf("invalid action type: %s", action)
	}

	// 1. Check if a wildcard already exists
	wildcardExists, err := s.permissionRepository.Check(ctx, userID, resourceType, a, nil)
	if err != nil {
		return err
	}

	if wildcardExists {
		// User already can do this action on ALL resources
		// No need to grant specific permission
		return nil
	}

	// 2. Check if specific permission already exists
	specificExists, err := s.permissionRepository.Check(ctx, userID, resourceType, a, resourceID)
	if err != nil {
		return err
	}

	if specificExists {
		// Already granted
		return nil
	}

	return s.permissionRepository.Create(ctx, &entities.Permission{
		UserID:       userID,
		Action:       a,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		GrantedAt:    time.Now(),
	})
}

func (s *PermissionService) GetByUser(ctx context.Context, userID string) (*[]PermissionSvc, error) {
	permissions, err := s.permissionRepository.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissionsSvc := make([]PermissionSvc, len(permissions))
	for i, p := range permissions {
		permissionsSvc[i] = PermissionSvc{
			ID:           p.ID,
			UserID:       p.UserID,
			Action:       p.Action,
			ResourceType: p.ResourceType,
			ResourceID:   p.ResourceID,
			GrantedAt:    p.GrantedAt,
			CreatedBy:    p.CreatedBy,
		}
	}

	return &permissionsSvc, nil
}

func (s *PermissionService) Check(ctx context.Context, userID, action, resourceType string, resourceID *string) (bool, error) {
	a := entities.Action(action)
	if !a.Valid() {
		return false, fmt.Errorf("invalid action type: %s", action)
	}

	allowed, err := s.permissionRepository.Check(ctx, userID, resourceType, a, resourceID)
	if err != nil {
		return false, err
	}

	return allowed, s.decisionRepository.Log(ctx, &entities.Decision{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Allowed:      allowed,
		Timestamp:    time.Now(), // the key to the log correction
	})
}

func (s *PermissionService) Delete(ctx context.Context, userID, action, resourceType string, resourceID *string) error {
	a := entities.Action(action)
	if !a.Valid() {
		return fmt.Errorf("invalid action type: %s", action)
	}

	// Find the exact permission
	perm, err := s.permissionRepository.Find(ctx, userID, a, resourceType, resourceID)
	if err != nil {
		return err
	}

	if perm == nil {
		return fmt.Errorf("permission not found")
	}

	// Delete by ID
	return s.permissionRepository.Delete(ctx, perm.ID)
}
