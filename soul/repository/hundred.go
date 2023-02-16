package repository

import (
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type HundredRepository struct {
	db *gorm.DB
}

func NewHundredRepository(db *gorm.DB) *HundredRepository {
	return &HundredRepository{db}
}

func (hr *HundredRepository) SaveHundred(rank uint, hundred uint) (model.Hundred, error) {
	var existing model.Hundred
	err := hr.db.Where("rank = ?", rank).Where("hundred = ?", hundred).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Hundred{}, err
	} else if err == nil {
		return existing, nil
	}
	h := model.Hundred{
		Rank:    rank,
		Hundred: hundred,
	}
	err = hr.db.Save(&h).Error
	return h, err
}

func (hr *HundredRepository) GetHundred(rank uint, hundred uint) (model.Hundred, error) {
	var h model.Hundred
	err := hr.db.Where("rank = ?", rank).Where("hundred = ?", hundred).First(&h).Error
	return h, err
}
