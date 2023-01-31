package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db}
}

func (pr *CardRepository) GetLastCards(ctx context.Context, count uint) ([]model.Card, error) {
	cards := make([]model.Card, 0, count)
	err := pr.db.
		Joins("Spell").
		Limit(int(count)).
		Order("cards.id desc").
		Limit(int(count)).
		Find(&cards).
		Error
	if err != nil {
		return []model.Card{}, errors.Wrapf(err, "failed to get cards count=%d", count)
	}

	return cards, err
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
	image := model.Image{}
	err := pr.db.Where("card_id = ?", cardID).First(&image).Error
	return image, err
}
