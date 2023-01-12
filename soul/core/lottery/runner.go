package lottery

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type lotteryRepository interface {
	SaveLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error)
}

type cardsRepository interface {
	GetCardsIDsByPeriod(ctx context.Context, start time.Time, end time.Time) ([]uint64, error)
}

type origin interface {
	Select(ctx context.Context, totalVariants uint64, saveDecision bool) (uint64, error)
}

type Runner struct {
	lotteryRepository lotteryRepository
	cardsRepository   cardsRepository
	origin            origin
}

func NewRunner(lotteryRepository lotteryRepository, cardsRepository cardsRepository, origin origin) *Runner {
	return &Runner{lotteryRepository, cardsRepository, origin}
}

func (lr *Runner) RunLottery(ctx context.Context, lottery model.Lottery) error {
	log.Info().Msgf("[runner] Running lottery: id=%d", lottery.ID)

	if lottery.State == model.LotteryStateFinished {
		return errors.Errorf("[runner] cannot run lottery(id=%d) in finished state", lottery.ID)
	}

	if lottery.State == model.LotteryStateWaiting {
		var err error
		if lottery, err = lr.startLottery(ctx, lottery); err != nil {
			return errors.Wrap(err, "failed to init lottery")
		}
		return nil // in this step only activate lottery
	}

	if lottery, finished, err := lr.performLotteryStep(ctx, lottery); err != nil {
		return errors.Wrapf(err, "[runner] failed performLotteryStep on lottery(id=%d)", lottery.ID)
	} else if finished {
		lottery, err := lr.finishLottery(ctx, lottery)
		if err != nil {
			return errors.Wrapf(err, "[runner] failed to finishLottery(id=%d)", lottery.ID)
		} else {
			log.Info().Msgf("[runner] lottery %d finished", lottery.ID)
		}
	}

	return nil
}

func (lr *Runner) startLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error) {
	log.Info().Msgf("[lottery_runner] startLottery id=%d", lottery.ID)
	lottery.State = model.LotteryStateRunning
	lottery.WinnersJSON = "[]"
	lottery.Started = time.Now()

	// select from 10 to 100 winners
	totalWinners, err := lr.origin.Select(ctx, 90, false)
	if err != nil {
		return model.Lottery{}, errors.Wrapf(err, "[runner] failed to get total winners from origin (lottery=%d)", lottery.ID)
	}
	totalWinners += 10 // min 10 winners, max 110

	lottery.TotalWinners = totalWinners
	lottery, err = lr.lotteryRepository.SaveLottery(ctx, lottery)
	if err != nil {
		return model.Lottery{}, errors.Wrapf(err, "[runner] failed to save lottery(id=%d)", lottery.ID)
	}
	log.Info().Msgf("[runner] started lottery(id=%d) with total winners %d", lottery.ID, lottery.TotalWinners)
	return lottery, nil
}

func (lr *Runner) performLotteryStep(ctx context.Context, lottery model.Lottery) (model.Lottery, bool, error) {
	var winners []uint64
	if err := json.Unmarshal([]byte(lottery.WinnersJSON), &winners); err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "[runner] failed to unmarshal lottery winners")
	}
	if uint64(len(winners)) >= lottery.TotalWinners {
		// need to finish lottery
		return lottery, true, nil
	}
	cards, err := lr.cardsRepository.GetCardsIDsByPeriod(ctx, lottery.CollectPeriodStart, lottery.CollectPeriodEnd)
	if err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "[runner] failed to GetCardsIDsByPeriod for lottery(id=%d)", lottery.ID)
	}

	totalCards := uint64(len(cards))
	selection, err := lr.origin.Select(ctx, totalCards-1, false)
	if err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "[runner] failed to select from origin with max=%d", totalCards)
	}
	log.Info().Msgf("[runner] selected %d(card id=%d), total cards %d", selection, cards[selection], len(cards))

	winners = append(winners, cards[selection])
	data, err := json.Marshal(winners)
	if err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "failed to marshal winners")
	}
	lottery.WinnersJSON = string(data)
	lottery, err = lr.lotteryRepository.SaveLottery(ctx, lottery)

	return lottery, false, err
}

func (lr *Runner) finishLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error) {
	lottery.State = model.LotteryStateFinished
	lottery.Finished = time.Now()
	return lr.lotteryRepository.SaveLottery(ctx, lottery)
}
