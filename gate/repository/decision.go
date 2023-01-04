package repository

import (
	"context"
	"github.com/artchitector/artchitect.git/gate/model"
	"gorm.io/gorm"
)

type DecisionRepository struct {
	db *gorm.DB
}

func NewDecisionRepository(db *gorm.DB) *DecisionRepository {
	return &DecisionRepository{db}
}

func (dr *DecisionRepository) GetLastDecision(ctx context.Context) (model.Decision, error) {
	decision := model.Decision{}
	err := dr.db.Last(&decision).Error
	return decision, err
}
