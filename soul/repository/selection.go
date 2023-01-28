package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
)

type SelectionRepository struct {
	db *gorm.DB
}

func NewSelectionRepository(db *gorm.DB) *SelectionRepository {
	return &SelectionRepository{db}
}

func (r *SelectionRepository) SaveSelection(ctx context.Context, selected model.Selection) (model.Selection, error) {
	err := r.db.Save(&selected).Error
	return selected, err
}
