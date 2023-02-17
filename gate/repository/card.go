package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

func (pr *CardRepository) GetCardsByRange(start uint, end uint) ([]model.Card, error) {
	var cards []model.Card
	log.Info().Msgf("[card_repo] get cards between %d and %d", start, end)
	err := pr.db.Joins("Spell").Where("cards.id between ? and ?", start, end).Find(&cards).Error
	return cards, err
}

func (pr *CardRepository) GetCards(ctx context.Context, IDs []uint) ([]model.Card, error) {
	var cards []model.Card
	err := pr.db.Joins("Spell").Where("cards.id in (?)", IDs).Find(&cards).Error
	return cards, err
}
