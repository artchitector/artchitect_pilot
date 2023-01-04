package repository

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/model"
	"gorm.io/gorm"
)

type StateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) *StateRepository {
	return &StateRepository{db}
}

func (dr *StateRepository) SaveState(ctx context.Context, state model.State) (model.State, error) {
	err := dr.db.Save(&state).Error
	return state, err
}
