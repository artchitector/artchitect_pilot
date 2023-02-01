package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
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

func (pr *CardRepository) DeleteCard(ctx context.Context, cardID uint) error {
	err := pr.db.Where("id = ?", cardID).Delete(&model.Card{}).Error
	return err
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

func (pr *CardRepository) GetCard(ctx context.Context, ID uint) (model.Card, error) {
	card := model.Card{}
	err := pr.db.
		Joins("Spell").
		Where("cards.id = ?", ID).
		Last(&card).
		Error
	if err != nil {
		return card, errors.Wrapf(err, "[card_repository] failed to find card %d", ID)
	} else {
		return card, nil
	}
}

func (pr *CardRepository) GetImage(ctx context.Context, cardID uint) (model.Image, error) {
	var image model.Image
	err := pr.db.Where("card_id=?", cardID).First(&image).Error
	return image, err
}
