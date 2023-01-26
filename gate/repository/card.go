package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type cache interface {
	RefreshLastCards(ctx context.Context, cards []model.Card) error
	SaveCard(ctx context.Context, card model.Card) error
}

type CardRepository struct {
	db    *gorm.DB
	cache cache
}

func NewCardRepository(db *gorm.DB, cache cache) *CardRepository {
	return &CardRepository{db, cache}
}

func (pr *CardRepository) GetLastCard(ctx context.Context) (model.Card, bool, error) {
	card := model.Card{}
	err := pr.db.Preload("Spell").Last(&card).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return card, false, nil
	} else if err != nil {
		return card, false, errors.Wrap(err, "failed to get last card")
	} else {
		return card, true, nil
	}
}

func (pr *CardRepository) GetLastCards(ctx context.Context, count uint) ([]model.Card, error) {
	cards := make([]model.Card, 0, count)
	err := pr.db.Preload("Spell").Limit(int(count)).Order("id desc").Limit(int(count)).Find(&cards).Error
	if err != nil {
		return []model.Card{}, errors.Wrapf(err, "failed to get cards count=%d", count)
	}
	go func() {
		if err := pr.cache.RefreshLastCards(ctx, cards); err != nil {
			log.Error().Err(err).Msgf("[card_repository] failed to reload last cards cache")
		}
	}()

	return cards, err
}

func (pr *CardRepository) GetCard(ctx context.Context, ID uint) (model.Card, bool, error) {
	card := model.Card{}
	err := pr.db.Preload("Spell").Where("id = ?", ID).Last(&card).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return card, false, nil
	} else if err != nil {
		return card, false, errors.Wrap(err, "failed to get last card")
	} else {
		go func() {
			if err := pr.cache.SaveCard(ctx, card); err != nil {
				log.Error().Err(err).Msgf("[card_repository] failed to save card(id=%d) to cache", card.ID)
			}
		}()
		return card, true, nil
	}
}
func (pr *CardRepository) GetCardsRange(ctx context.Context, from uint, to uint) ([]model.Card, error) {
	var min, max uint
	var order string
	if from < to {
		min = from
		max = to
		order = "id asc"
	} else {
		min = to
		max = from
		order = "id desc"
	}
	if max-min > 100 {
		return []model.Card{}, errors.Errorf("maximum 100 cards allowed")
	}
	cards := make([]model.Card, 0, max-min)
	err := pr.db.Preload("Spell").Where("id between ? and ?", min, max).Order(order).Find(&cards).Error
	return cards, err
}
