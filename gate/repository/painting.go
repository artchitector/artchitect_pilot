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

func (pr *CardRepository) GetLastCard(ctx context.Context) (model.Card, bool, error) {
	painting := model.Card{}
	err := pr.db.Preload("Spell").Last(&painting).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return painting, false, nil
	} else if err != nil {
		return painting, false, errors.Wrap(err, "failed to get last painting")
	} else {
		return painting, true, nil
	}
}

func (pr *CardRepository) GetLastCards(ctx context.Context, count uint64) ([]model.Card, error) {
	paintings := make([]model.Card, 0, count)
	err := pr.db.Preload("Spell").Limit(int(count)).Order("id desc").Find(&paintings).Error
	return paintings, err
}

func (pr *CardRepository) GetCard(ctx context.Context, ID uint) (model.Card, bool, error) {
	painting := model.Card{}
	err := pr.db.Preload("Spell").Where("id = ?", ID).Last(&painting).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return painting, false, nil
	} else if err != nil {
		return painting, false, errors.Wrap(err, "failed to get last painting")
	} else {
		return painting, true, nil
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
		return []model.Card{}, errors.Errorf("maximum 100 paintings allowed")
	}
	paintings := make([]model.Card, 0, max-min)
	err := pr.db.Preload("Spell").Where("id between ? and ?", min, max).Order(order).Find(&paintings).Error
	return paintings, err
}
