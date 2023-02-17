package repository

import (
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
)

type HundredRepository struct {
	db *gorm.DB
}

func NewHundredRepository(db *gorm.DB) *HundredRepository {
	return &HundredRepository{db}
}

// example return [10000/0, 10000/10000, 10000/20000]
func (hr *HundredRepository) FindAllTenK() ([]model.Hundred, error) {
	var hundreds []model.Hundred
	err := hr.db.Where("rank = ?", model.Rank10000).Order("hundred desc").Find(&hundreds).Error
	return hundreds, err
}

// example input 10000 means give some 1000 hundreds [10000, 11000, 12000...]
func (hr *HundredRepository) FindKList(tenKHundred uint) ([]model.Hundred, error) {
	var hundreds []model.Hundred
	err := hr.db.Where("rank = ?", model.Rank1000).Where("hundred between ? and ?", tenKHundred, tenKHundred+model.Rank10000-1).Order("hundred desc").Find(&hundreds).Error
	return hundreds, err
}

// example input 11000 means give some hundreds [11000, 11100, 11200...]
func (hr *HundredRepository) FindHList(kHundred uint) ([]model.Hundred, error) {
	var hundreds []model.Hundred
	err := hr.db.Where("rank = ?", model.Rank100).Where("hundred between ? and ?", kHundred, kHundred+model.Rank1000-1).Order("hundred desc").Find(&hundreds).Error
	return hundreds, err
}
