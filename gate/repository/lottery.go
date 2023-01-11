package repository

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
	"sort"
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

func (lr *LotteryRepository) GetSelection(ctx context.Context) ([]uint64, error) {
	selection := make(map[uint64]struct{})
	var lotteries []model.Lottery
	err := lr.db.Where("state = ?", model.LotteryStateFinished).Find(&lotteries).Error
	if err != nil {
		return []uint64{}, err
	}
	for _, lottery := range lotteries {
		var winners []uint64
		if err := json.Unmarshal([]byte(lottery.WinnersJSON), &winners); err != nil {
			return []uint64{}, err
		}
		for _, winner := range winners {
			selection[winner] = struct{}{}
		}
	}

	list := make([]uint64, 0, len(selection))
	for id, _ := range selection {
		list = append(list, id)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] > list[j]
	})
	return list, nil
}
