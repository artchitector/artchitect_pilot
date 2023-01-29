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

func (pr *CardRepository) SaveCard(ctx context.Context, painting model.Card) (model.Card, error) {
	err := pr.db.Save(&painting).Error
	return painting, err
}

func (pr *CardRepository) GetCardsIDsByPeriod(ctx context.Context, start time.Time, end time.Time) ([]uint, error) {
	var ids []uint
	err := pr.db.Model(&model.Card{}).Select("id").Where("created_at between ? and ?", start, end).Find(&ids).Error
	return ids, err
}

func (pr *CardRepository) GetTotalCards(ctx context.Context) (uint, error) {
	var count uint
	err := pr.db.Select("count(id)").Model(&model.Card{}).Find(&count).Error
	return count, err
}

func (pr *CardRepository) GetCardWithOffset(offset uint) (model.Card, error) {
	var card model.Card
	err := pr.db.
		Joins("Spell").
		Joins("Image").
		Order("cards.id asc").
		Limit(1).
		Offset(int(offset)).
		Find(&card).Error
	return card, err
}

func (pr *CardRepository) GetLastCardPaintTime(ctx context.Context) (uint, error) {
	var paintTime uint
	err := pr.db.Select("paint_time").Model(&model.Card{}).Order("id desc").Limit(1).Scan(&paintTime).Error
	return paintTime, err
}
