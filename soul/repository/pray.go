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

func (pr *PrayRepository) GetNextPray(ctx context.Context) (model.PrayWithQuestion, error) {
	var pray model.PrayWithQuestion
	err := pr.db.Where("state = ?", model.PrayStateWaiting).Order("id asc").Limit(1).First(&pray).Error
	return pray, err
}

func (pr *PrayRepository) AnswerPray(ctx context.Context, pray model.PrayWithQuestion, answer uint) error {
	pray.Answer = answer
	pray.State = model.PrayStateAnswered
	err := pr.db.Save(&pray).Error
	return err
}
