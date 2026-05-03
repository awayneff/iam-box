// app/service/decision.go

package service

import (
	"context"
	"iam-box/app/entities"
	"iam-box/app/repository"
	"time"
)

type Decision struct {
	ID           uint
	UserID       string
	Action       string
	ResourceType string
	ResourceID   *string
	Allowed      bool
	Reason       string
	Timestamp    time.Time
}

type DecisionService struct {
	decisionRepository repository.DecisionRepository
}

func NewDecisionService(decisionRepository repository.DecisionRepository) *DecisionService {
	return &DecisionService{
		decisionRepository: decisionRepository,
	}
}

func (s *DecisionService) Log(ctx context.Context, userID, action, resourceType, reason string, resourceID *string, allowed bool) error {
	return s.decisionRepository.Log(ctx, &entities.Decision{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Allowed:      allowed,
		Reason:       reason,
		Timestamp:    time.Now(),
	})
}

func (s *DecisionService) List(ctx context.Context, limit, offset int) ([]Decision, error) {
	d, err := s.decisionRepository.List(ctx, limit, offset)
	if err != nil {
		return []Decision{}, err
	}

	decisions := make([]Decision, len(d))
	for i, decision := range d {
		decisions[i] = Decision{
			ID:           decision.ID,
			UserID:       decision.UserID,
			Action:       decision.Action,
			ResourceType: decision.ResourceType,
			ResourceID:   decision.ResourceID,
			Allowed:      decision.Allowed,
			Reason:       decision.Reason,
			Timestamp:    decision.Timestamp,
		}
	}

	return decisions, nil
}
