package hundreds

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type hundredsRepo interface {
	GetHundred(rank uint, hundred uint) (model.Hundred, error)
}

type cardsRepo interface {
	GetMaxCardID(ctx context.Context) (uint, error)
}

type combinator interface {
	CombineHundred(ctx context.Context, rank uint, hundred uint) error
}

type HundredsWorker struct {
	cardsRepo         cardsRepo
	hundredsRepo      hundredsRepo
	combinator        combinator
	lastWorkedRank    uint
	lastWorkedHundred uint
}

func NewHundredsWorker(cardsRepo cardsRepo, hundredsRepo hundredsRepo, combinator combinator) *HundredsWorker {
	return &HundredsWorker{cardsRepo, hundredsRepo, combinator, 0, 0}
}

func (w *HundredsWorker) Work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.NewTicker(time.Second).C:
			allDone, err := w.WorkOnce(ctx)
			if err != nil {
				log.Error().Err(err).Send()
			}
			if allDone {
				log.Info().Msgf("[hundreds_worker] finished work! all done")
				return
			}
		}
	}
}

func (w *HundredsWorker) WorkOnce(ctx context.Context) (bool, error) {
	maxCardID, err := w.cardsRepo.GetMaxCardID(ctx)
	if err != nil {
		return false, err
	}
	ranks := []uint{model.Rank10000, model.Rank1000, model.Rank100}
	for _, r := range ranks {
		if w.lastWorkedRank < r {
			continue
		}
		for h := w.lastWorkedHundred; h < maxCardID; h += r {
			log.Info().Msgf("[hundreds_worker] Starting work on r:%d h:%d", r, h)
			_, err := w.hundredsRepo.GetHundred(r, h)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return false, errors.Wrapf(err, "[hundreds_worker] failed to get hundred")
			} else if err == nil {
				// hundred already exists
				continue
			}
			err = w.combinator.CombineHundred(ctx, r, h)
			if err != nil {
				return false, errors.Wrapf(err, "[hundreds_worker] failed to combine r:%d h:%d", r, h)
			}
		}
	}
	return true, nil
}
