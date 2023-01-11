package lottery

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type lotteryRepository interface {
	StartLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error)
	SaveLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error)
	SaveTourWinners(ctx context.Context, tour model.LotteryTour, winners []uint64) (model.LotteryTour, error)
	FinishTour(ctx context.Context, tour model.LotteryTour) (model.LotteryTour, error)
	GetNextAvailableTour(ctx context.Context, lotteryID uint) (model.LotteryTour, error)
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
		if lottery, err = lr.initLottery(ctx, lottery); err != nil {
			return errors.Wrap(err, "failed to init lottery")
		}
	}

	if len(lottery.Tours) == 0 {
		return errors.Errorf("failed to start lottery. no tours (id=%d)", lottery.ID)
	}

	// run tours
	nextTour, found, err := lr.getNextTour(ctx, lottery)
	if err != nil {
		return errors.Wrapf(err, "failed to get next tour for lottery(id=%d)", lottery.ID)
	}
	if !found {
		// no next tour
		_, err := lr.finishLottery(ctx, lottery)
		return err
	} else {
		err := lr.runTour(ctx, lottery, nextTour)
		return err
	}
}

func (lr *Runner) initLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error) {
	log.Info().Msgf("[lottery_runner] initLottery id=%d", lottery.ID)
	return lr.lotteryRepository.StartLottery(ctx, lottery)
}

func (lr *Runner) getNextTour(ctx context.Context, lottery model.Lottery) (model.LotteryTour, bool, error) {
	tour, err := lr.lotteryRepository.GetNextAvailableTour(ctx, lottery.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.LotteryTour{}, false, nil
	} else if err != nil {
		return model.LotteryTour{}, false, err
	}
	return tour, true, nil
}

func (lr *Runner) finishLottery(ctx context.Context, lottery model.Lottery) (model.Lottery, error) {
	log.Info().Msgf("finish lottery id=%d", lottery.ID)
	lastTour := lottery.Tours[len(lottery.Tours)-1]
	lottery.WinnersJSON = lastTour.WinnersJSON
	lottery.State = model.LotteryStateFinished
	var err error
	lottery, err = lr.lotteryRepository.SaveLottery(ctx, lottery)
	if err != nil {
		return model.Lottery{}, errors.Wrapf(err, "failed to finish lottery")
	}
	log.Info().Msgf("[runner] finished lottery(id=%d). winners: %s", lottery.ID, lottery.WinnersJSON)
	return lottery, nil
}

func (lr *Runner) runTour(ctx context.Context, lottery model.Lottery, tour model.LotteryTour) error {
	// each tour have MaxWinners, number of resulting cards. Need to select this resulting cards and take best from them
	var currentWinners []uint64
	if err := json.Unmarshal([]byte(tour.WinnersJSON), &currentWinners); err != nil {
		return errors.Wrapf(err, "failed unmarshal winners json")
	}
	if uint64(len(currentWinners)) == tour.MaxWinners {
		// tour finished
		if _, err := lr.lotteryRepository.FinishTour(ctx, tour); err != nil {
			return errors.Wrapf(err, "failed finish tour")
		} else {
			log.Info().Msgf("[runner] tour finished id=%d", tour.ID)
			return nil
		}
	}

	// if it's first tour, then we take all cards from period
	// if it's second tour, then take winners of previous tour
	var cards []uint64

	if lottery.IsFirstTour(tour) {
		var err error
		cards, err = lr.cardsRepository.GetCardsIDsByPeriod(ctx, lottery.CollectPeriodStart, lottery.CollectPeriodEnd)
		if errors.Is(err, gorm.ErrRecordNotFound) || len(cards) == 0 {
			// no cards to select
			return errors.Errorf("no cards found for lottery(id=%d), tour(id=%d)", lottery.ID, tour.ID)
		} else if err != nil {
			return errors.Wrapf(err, "failed get for lottery(id=%d), tour(id=%d)", lottery.ID, tour.ID)
		}
		log.Info().Msgf("[runner] first tour(id=%d) of lottery(id=%d) started. got %d cards", tour.ID, lottery.ID, len(cards))
	} else {
		prevTour, err := lottery.PreviousTour(tour)
		if err != nil {
			return err
		}
		winnersJSON := prevTour.WinnersJSON
		var previousWinners []uint64
		if err := json.Unmarshal([]byte(winnersJSON), &previousWinners); err != nil {
			return errors.Wrapf(err, "failed to parse previous tour winners")
		}
		cards = previousWinners
		log.Info().Msgf("[runner] next tour(id=%d) of lottery(id=%d) started. got previous winners: %+v", tour.ID, lottery.ID, cards)
	}

	totalCardsLength := uint64(len(cards))
	log.Info().Msgf("[runner] %d cards found to tour(id=%d)", totalCardsLength, tour.ID)

	selection, err := lr.origin.Select(ctx, totalCardsLength-1, false)
	if err != nil {
		return errors.Wrapf(err, "failed to get origin for card selection, tour(id=%d)", tour.ID)
	}
	log.Info().Msgf("[runner] got selection %d from origin. it is card(id=%d)", selection, cards[selection])

	currentWinners = append(currentWinners, cards[selection])
	tour, err = lr.lotteryRepository.SaveTourWinners(ctx, tour, currentWinners)
	if err != nil {
		return errors.Wrapf(err, "failed to set winners in tour(id=%d)", tour.ID)
	}
	log.Info().Msgf("[runner] saved winners tour(id=%d): %s", tour.ID, tour.WinnersJSON)

	return nil
}
