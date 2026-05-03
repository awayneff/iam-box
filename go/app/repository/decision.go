// app/repository/decision.go

package repository

import (
	"context"
	"iam-box/app/entities"

	"gorm.io/gorm"
)

type DecisionRepository struct {
	db *gorm.DB
}

func NewDecisionRepository(db *gorm.DB) *DecisionRepository {
	return &DecisionRepository{db: db}
}

func (r *DecisionRepository) Log(ctx context.Context, d *entities.Decision) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *DecisionRepository) List(ctx context.Context, limit, offset int) ([]entities.Decision, error) {
	var decisions []entities.Decision
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("timestamp DESC").Find(&decisions).Error
	return decisions, err
}
