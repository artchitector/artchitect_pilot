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

type selectionRepository interface {
	SaveSelection(ctx context.Context, selected model.Selection) (model.Selection, error)
}

type cardsRepository interface {
	GetCardsIDsByPeriod(ctx context.Context, start time.Time, end time.Time) ([]uint, error)
}

type notifier interface {
	NotifyNewSelection(ctx context.Context, selection model.Selection) error
	NotifyLottery(ctx context.Context, state model.LotteryState) error
}

type origin interface {
	Select(ctx context.Context, totalVariants uint) (uint, error)
}

type Runner struct {
	lotteryRepository   lotteryRepository
	selectionRepository selectionRepository
	cardsRepository     cardsRepository
	origin              origin
	notifier            notifier
}

func NewRunner(
	lotteryRepository lotteryRepository,
	selectionRepository selectionRepository,
	cardsRepository cardsRepository,
	origin origin,
	notifier notifier,
) *Runner {
	return &Runner{
		lotteryRepository,
		selectionRepository,
		cardsRepository,
		origin,
		notifier,
	}
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
	totalWinners, err := lr.origin.Select(ctx, 90)
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
	var winners []uint
	if err := json.Unmarshal([]byte(lottery.WinnersJSON), &winners); err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "[runner] failed to unmarshal lottery winners")
	}
	if uint(len(winners)) >= lottery.TotalWinners {
		// need to finish lottery
		return lottery, true, nil
	}
	cards, err := lr.cardsRepository.GetCardsIDsByPeriod(ctx, lottery.CollectPeriodStart, lottery.CollectPeriodEnd)
	if err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "[runner] failed to GetCardsIDsByPeriod for lottery(id=%d)", lottery.ID)
	}

	totalCards := uint(len(cards))
	selection, err := lr.origin.Select(ctx, totalCards)
	if err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "[runner] failed to select from origin with max=%d", totalCards)
	}
	cardID := cards[selection]
	log.Info().Msgf("[runner] selected %d(card id=%d), total cards %d", selection, cardID, len(cards))

	winners = append(winners, cardID)
	data, err := json.Marshal(winners)
	if err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "failed to marshal winners")
	}
	lottery.WinnersJSON = string(data)
	lottery, err = lr.lotteryRepository.SaveLottery(ctx, lottery)
	if err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "[runner] failed to save lottery winners (id=%d)", lottery.ID)
	}
	lr.notifyLottery(ctx, lottery, 0, 0)

	selected, err := lr.selectionRepository.SaveSelection(ctx, model.Selection{
		CardID:    cardID,
		LotteryID: lottery.ID,
	})
	if err != nil {
		return model.Lottery{}, false, errors.Wrapf(err, "[runner] failed to save selected card (lottery=%d, card=%d)", lottery.ID, cardID)
	}
	log.Info().Msgf("[runner] saved selected (lottery=%d, card=%d)", selected.LotteryID, selected.CardID)
	if err := lr.notifier.NotifyNewSelection(ctx, selected); err != nil {
		log.Error().Err(err).Msgf("[runner] failed to notify selection (id=%d)", selected.ID)
	}

	return lottery, false, err
}

func (lr *Runner) finishLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error) {
	lottery.State = model.LotteryStateFinished
	lottery.Finished = time.Now()
	saved, err := lr.lotteryRepository.SaveLottery(ctx, lottery)
	if err != nil {
		return model.Lottery{}, errors.Wrapf(err, "[runner] failed to save lottery id=%d", lottery.ID)
	}
	start := time.Now()
	var enjoySeconds uint = 10
	enjoyFinish := start.Add(time.Second * time.Duration(enjoySeconds))
	lr.notifyLottery(ctx, saved, 0, enjoySeconds)
forLoop:
	for {
		select {
		case <-ctx.Done():
			break forLoop
		case <-time.Tick(time.Second):
			if time.Now().After(enjoyFinish) {
				lr.notifyLottery(ctx, saved, enjoySeconds, enjoySeconds)
				break forLoop
			} else {
				currentEnjoySeconds := uint(time.Now().Sub(start).Seconds())
				lr.notifyLottery(ctx, saved, currentEnjoySeconds, enjoySeconds)
			}
		}
	}
	return saved, nil
}

func (lr *Runner) notifyLottery(ctx context.Context, lottery model.Lottery, enjoyTime uint, totalEnjoyTime uint) {
	var winners []uint
	if err := json.Unmarshal([]byte(lottery.WinnersJSON), &winners); err != nil {
		log.Error().Err(err).Msgf("[runner] failed to unmarshal lottery winner (id=%d). Json: %s", lottery.ID, lottery.WinnersJSON)
	}
	lottery.Winners = winners
	state := model.LotteryState{
		Lottery:          lottery,
		EnjoyCurrentTime: enjoyTime,
		EnjoyTotalTime:   totalEnjoyTime,
	}
	if err := lr.notifier.NotifyLottery(ctx, state); err != nil {
		log.Error().Err(err).Msgf("[runner] failed to notify lottery (id=%d)", lottery.ID)
	}
}
