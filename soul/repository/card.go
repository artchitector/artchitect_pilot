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

func (pr *CardRepository) GetTotalCards(ctx context.Context) (uint64, error) {
	var count uint64
	err := pr.db.Select("count(id)").Model(&model.Card{}).Find(&count).Error
	return count, err
}

func (pr *CardRepository) GetCardWithOffset(offset uint64) (model.Card, error) {
	var card model.Card
	err := pr.db.Preload("Spell").Order("id asc").Limit(1).Offset(int(offset)).Find(&card).Error
	return card, err
}
