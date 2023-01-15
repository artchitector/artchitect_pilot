package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
)

type PrayRepository struct {
	db *gorm.DB
}

func NewPrayRepository(db *gorm.DB) *PrayRepository {
	return &PrayRepository{db: db}
}

func (pr *PrayRepository) MakePray(ctx context.Context) (model.PrayWithQuestion, error) {
	pray := model.PrayWithQuestion{
		Model:  gorm.Model{},
		State:  model.PrayStateWaiting,
		Answer: 0,
	}
	err := pr.db.Create(&pray).Error
	return pray, err
}

func (pr *PrayRepository) GetAnswer(ctx context.Context, prayId uint64) (uint64, error) {
	var answer uint64
	err := pr.db.Model(&model.PrayWithQuestion{}).Select("answer").Where("id = ?", prayId).Where("state = ?", model.PrayStateAnswered).Scan(&answer).Error
	return answer, err
}
