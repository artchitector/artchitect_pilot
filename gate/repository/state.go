package repository

import (
	"context"
	"github.com/artchitector/artchitect.git/gate/model"
	"gorm.io/gorm"
)

type StateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) *StateRepository {
	return &StateRepository{db}
}

func (dr *StateRepository) GetLastState(ctx context.Context) (model.State, error) {
	state := model.State{}
	err := dr.db.Last(&state).Error
	return state, err
}
