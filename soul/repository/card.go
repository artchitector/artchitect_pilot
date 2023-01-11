package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
	"time"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db}
}

func (pr *CardRepository) SavePainting(ctx context.Context, painting model.Card) (model.Card, error) {
	err := pr.db.Save(&painting).Error
	return painting, err
}

func (pr *CardRepository) GetCardsIDsByPeriod(ctx context.Context, start time.Time, end time.Time) ([]uint64, error) {
	var ids []uint64
	err := pr.db.Model(&model.Card{}).Select("id").Where("created_at between ? and ?", start, end).Find(&ids).Error
	return ids, err
}
