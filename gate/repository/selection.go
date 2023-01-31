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

func (r *SelectionRepository) GetSelection(ctx context.Context) ([]uint, error) {
	var ids []uint
	err := r.db.Select("distinct(card_id)").Model(&model.Selection{}).Order("card_id DESC").Find(&ids).Error
	return ids, err
}

func (r *SelectionRepository) GetSelectionLimit(ctx context.Context, limit int) ([]uint, error) {
	var ids []uint
	err := r.db.
		Select("distinct(card_id)").
		Model(&model.Selection{}).
		Order("card_id DESC").
		Limit(limit).
		Find(&ids).Error
	return ids, err
}
