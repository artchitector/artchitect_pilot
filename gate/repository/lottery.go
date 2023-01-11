package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
)

type LotteryRepository struct {
	db *gorm.DB
}

func NewLotteryRepository(db *gorm.DB) *LotteryRepository {
	return &LotteryRepository{db}
}

func (lr *LotteryRepository) GetActiveLottery(ctx context.Context) (model.Lottery, error) {
	var lottery model.Lottery
	err := lr.db.
		Preload("Tours").
		Where("(state = ? or state = ?)", model.LotteryStateRunning, model.LotteryStateFinished).
		Where("start_time < current_timestamp").
		First(&lottery).
		Error
	return lottery, err
}

func (lr *LotteryRepository) GetLastLotteries(ctx context.Context, lastN uint) ([]model.Lottery, error) {
	var lotteries []model.Lottery
	err := lr.db.
		Preload("Tours").
		Order("id desc").Limit(int(lastN)).Find(&lotteries).
		Error
	return lotteries, err
}
