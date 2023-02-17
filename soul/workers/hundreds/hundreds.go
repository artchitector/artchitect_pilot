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
	lastWorkedHundred map[uint]uint
}

func NewHundredsWorker(cardsRepo cardsRepo, hundredsRepo hundredsRepo, combinator combinator) *HundredsWorker {
	return &HundredsWorker{cardsRepo, hundredsRepo, combinator, make(map[uint]uint)}
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
	/*
		Этот воркер должен создавать карточки сборных сотен, тысяч и десятитысяч, которые не были созданы ранее
		Он сначала проходит по всем десятитысячам, чтобы они были созданы, затем по всем тысячам, затем по всем сотням, и постепенно их создаёт
	*/
	maxCardID, err := w.cardsRepo.GetMaxCardID(ctx)
	if err != nil {
		return false, err
	}
	maxCardID -= maxCardID % model.Rank100
	log.Info().Msgf("[hundreds_worker] used maxCardID=%d", maxCardID)

	allDone, err := w.WorkOnceWithRank(ctx, model.Rank10000, maxCardID)
	if err != nil {
		return false, errors.Wrapf(err, "[hundreds_worker] failed to work with rank %d", model.Rank10000)
	}
	if !allDone {
		return false, nil
	}

	allDone, err = w.WorkOnceWithRank(ctx, model.Rank1000, maxCardID)
	if err != nil {
		return false, errors.Wrapf(err, "[hundreds_worker] failed to work with rank %d", model.Rank1000)
	}
	if !allDone {
		return false, nil
	}

	allDone, err = w.WorkOnceWithRank(ctx, model.Rank100, maxCardID)
	if err != nil {
		return false, errors.Wrapf(err, "[hundreds_worker] failed to work with rank %d", model.Rank100)
	}
	if !allDone {
		return false, nil
	}

	return true, nil
}

func (w *HundredsWorker) WorkOnceWithRank(ctx context.Context, rank uint, maxCardID uint) (bool, error) {
	start := uint(0)
	if last, ok := w.lastWorkedHundred[rank]; ok {
		start = last + rank
	}
	for i := start; i < maxCardID; i += rank {
		if rank == model.Rank100 && i == 0 {
			// empty 100, skip
			w.lastWorkedHundred[rank] = i
			continue
		}
		_, err := w.hundredsRepo.GetHundred(rank, i)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.Wrapf(err, "[hundreds_worker] failed to find existing hundred")
		} else if err == nil {
			log.Info().Msgf("[hundreds_worker] hundred r:%d h:%d already exists", rank, i)
			// hundred already exists
			w.lastWorkedHundred[rank] = i
			continue
		}
		log.Info().Msgf("[hundreds_worker] combining r:%d h:%d", rank, i)
		err = w.combinator.CombineHundred(ctx, rank, i)
		if err != nil {
			return false, errors.Wrapf(err, "[hundreds_worker] failed to combine r:%d h:%d", rank, i)
		}
		log.Info().Msgf("[hundreds_worker] combined r:%d h:%d", rank, i)
		w.lastWorkedHundred[rank] = i
		return false, nil
	}
	return true, nil
}
