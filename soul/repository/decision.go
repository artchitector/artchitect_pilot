package repository

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/model"
	"gorm.io/gorm"
)

type DecisionRepository struct {
	db *gorm.DB
}

func NewDecisionRepository(db *gorm.DB) *DecisionRepository {
	return &DecisionRepository{db}
}

func (dr *DecisionRepository) SaveDecision(ctx context.Context, decision model.Decision) (model.Decision, error) {
	err := dr.db.Save(&decision).Error
	return decision, err
}
