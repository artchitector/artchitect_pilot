package repository

import (
	"context"
	"encoding/json"
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
		Preload("Tours").
		Where("(state = ? or state = ?)", model.LotteryStateWaiting, model.LotteryStateRunning).
		Where("start_time < current_timestamp").
		First(&lottery).
		Error
	return lottery, err
}

func (lr *LotteryRepository) GetLottery(ctx context.Context, ID uint64) (model.Lottery, error) {
	var lottery model.Lottery
	err := lr.db.
		Preload("Tours").
		Where("id = ?", ID).
		First(&lottery).
		Error
	return lottery, err
}

func (lr *LotteryRepository) StartLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error) {
	//change lottery status
	lottery.State = model.LotteryStateRunning
	lottery.WinnersJSON = "[]"
	if err := lr.db.Save(&lottery).Error; err != nil {
		return lottery, err
	}
	// reset tours
	if err := lr.db.
		Model(&model.LotteryTour{}).
		Where("lottery_id = ?", lottery.ID).
		Update("state", model.LotteryStateWaiting).
		Error; err != nil {
		return lottery, errors.Wrap(err, "failed to reset tours states")
	}

	return lr.GetLottery(ctx, uint64(lottery.ID))
}

func (lr *LotteryRepository) SaveTourWinners(ctx context.Context, tour model.LotteryTour, winners []uint64) (model.LotteryTour, error) {
	bytes, err := json.Marshal(winners)
	if err != nil {
		return model.LotteryTour{}, err
	}
	tour.WinnersJSON = string(bytes)
	err = lr.db.Save(&tour).Error
	return tour, err
}

func (lr *LotteryRepository) FinishTour(ctx context.Context, tour model.LotteryTour) (model.LotteryTour, error) {
	tour.State = model.LotteryStateFinished
	err := lr.db.Save(&tour).Error
	return tour, err
}

func (lr *LotteryRepository) SaveLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error) {
	err := lr.db.Save(&lottery).Error
	return lottery, err
}

func (lr *LotteryRepository) GetNextAvailableTour(ctx context.Context, lotteryID uint) (model.LotteryTour, error) {
	var tour model.LotteryTour
	err := lr.db.Where("lottery_id = ? and state in ?", lotteryID, []string{model.LotteryStateWaiting, model.LotteryStateRunning}).Order("id asc").First(&tour).Error
	return tour, err
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

	tours := []model.LotteryTour{
		{
			Name:        "most100",
			MaxWinners:  100,
			WinnersJSON: "[]",
			State:       model.LotteryStateWaiting,
		},
		{
			Name:        "top10",
			MaxWinners:  10,
			WinnersJSON: "[]",
			State:       model.LotteryStateWaiting,
		},
	}
	lottery := model.Lottery{
		Name:               fmt.Sprintf("daily lottery %s", collectPeriodStart.Format("2 Jan 2006")),
		Type:               model.LotteryTypeDaily,
		StartTime:          lotteryStartTime,
		CollectPeriodStart: collectPeriodStart,
		CollectPeriodEnd:   collectPeriodEnd,
		State:              model.LotteryStateWaiting,
		TotalWinners:       10,
		Tours:              tours,
		WinnersJSON:        "[]",
	}

	if err := lr.db.Save(&lottery).Error; err != nil {
		return err
	} else {
		log.Info().Msgf("[lottery] generated new daily lottery (id=%d)", lottery.ID)
		return nil
	}
}

func (lr *LotteryRepository) GetLastTour(ctx context.Context, lottery model.Lottery) (model.LotteryTour, error) {
	var tour model.LotteryTour
	err := lr.db.Where("lottery_id = ?", lottery.ID).Order("id desc").First(&tour).Error
	return tour, err
}
