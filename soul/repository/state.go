package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type StateRepository struct {
	db *gorm.DB
}

func NewStateRepository(db *gorm.DB) *StateRepository {
	return &StateRepository{db}
}

func (dr *StateRepository) SaveState(ctx context.Context, state model.State) (model.State, error) {
	var lastState model.State
	if err := dr.db.Last(&lastState).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		lastState = model.State{State: model.StateNotWorking}
	} else if err != nil {
		return model.State{}, errors.Wrap(err, "failed to get last state")
	}
	if lastState.State == state.State {
		return lastState, nil
	}
	err := dr.db.Save(&state).Error
	return state, err
}
