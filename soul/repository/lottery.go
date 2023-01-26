package repository

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
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
		Where("(state = ? or state = ?)", model.LotteryStateWaiting, model.LotteryStateRunning).
		Where("start_time < current_timestamp").
		First(&lottery).
		Error
	return lottery, err
}

func (lr *LotteryRepository) GetLottery(ctx context.Context, ID uint) (model.Lottery, error) {
	var lottery model.Lottery
	err := lr.db.
		Where("id = ?", ID).
		First(&lottery).
		Error
	return lottery, err
}

func (lr *LotteryRepository) SaveLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error) {
	err := lr.db.Save(&lottery).Error
	return lottery, err
}

func (lr *LotteryRepository) InitDailyLottery(ctx context.Context) error {
	mow, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return errors.Wrap(err, "failed to get Europe/Moscow tz")
	}

	today := time.Now().Add(time.Hour)
	collectPeriodStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, mow)
	collectPeriodEnd := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, mow)

	var existingDailyLottery model.Lottery
	err = lr.db.
		Where("state = ? and start_time > ? and type = ?", model.LotteryStateWaiting, collectPeriodEnd, model.LotteryTypeDaily).
		First(&existingDailyLottery).
		Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else if err == nil {
		log.Info().Msgf("daily lottery already set (id=%d)", existingDailyLottery.ID)
		return nil
	}

	lotteryStartTime := collectPeriodEnd.Add(time.Second * 2)

	lottery := model.Lottery{
		Name:               fmt.Sprintf("daily lottery %s", collectPeriodStart.Format("2 Jan 2006")),
		Type:               model.LotteryTypeDaily,
		StartTime:          lotteryStartTime,
		CollectPeriodStart: collectPeriodStart,
		CollectPeriodEnd:   collectPeriodEnd,
		State:              model.LotteryStateWaiting,
		TotalWinners:       0,
		WinnersJSON:        "[]",
	}

	if err := lr.db.Save(&lottery).Error; err != nil {
		return err
	} else {
		log.Info().Msgf("[lottery] generated new daily lottery (id=%d)", lottery.ID)
		return nil
	}
}
