package artchitect

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type notifier interface {
	NotifyTick(ctx context.Context, tick int) error
	NotifyNewCard(ctx context.Context, card model.Card) error
}

type ArtistContract interface {
	GetCard(ctx context.Context, spell model.Spell) (model.Card, error)
}

type state interface {
	SetState(ctx context.Context, state string)
}

type speller interface {
	MakeSpell(ctx context.Context) (model.Spell, error)
}

type merciful interface {
	AnswerPray(ctx context.Context) (bool, error)
}

type lotteryRepository interface {
	GetActiveLottery(ctx context.Context) (model.Lottery, error)
	InitDailyLottery(ctx context.Context) error
}

type lotteryRunner interface {
	RunLottery(ctx context.Context, lottery model.Lottery) error
}

type Config struct {
	CardsCreationEnabled bool
	LotteryEnabled       bool
	MercifulEnabled      bool
}

type Artchitect struct {
	config            Config
	state             state
	speller           speller
	artist            ArtistContract
	lotteryRepository lotteryRepository
	lotteryRunner     lotteryRunner
	merciful          merciful
	notifier          notifier
}

func NewArtchitect(
	config Config,
	state state,
	speller speller,
	artist ArtistContract,
	lotteryRepository lotteryRepository,
	lotteryRunner lotteryRunner,
	merciful merciful,
	notifier notifier,
) *Artchitect {
	return &Artchitect{
		config,
		state,
		speller,
		artist,
		lotteryRepository,
		lotteryRunner,
		merciful,
		notifier,
	}
}

func (a *Artchitect) Run(ctx context.Context, tick int) error {
	log.Info().Msgf("[artchitect] tick=%d", tick)
	if err := a.notifier.NotifyTick(ctx, tick); err != nil {
		log.Error().Err(err).Send()
	}
	if tick%10 == 0 {
		return a.maintenance(ctx)
	}
	if a.config.LotteryEnabled {
		activeLottery, err := a.lotteryRepository.GetActiveLottery(ctx)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "failed to get active lottery")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return a.runLottery(ctx, activeLottery)
		}
	}
	if a.config.MercifulEnabled {
		answered, err := a.merciful.AnswerPray(ctx)
		if err != nil {
			return errors.Wrap(err, "[artchitect] failed pray answer")
		} else if answered {
			log.Info().Msgf("[artchitect] answered a pray")
			return nil
		}
	}
	if a.config.CardsCreationEnabled {
		return a.runCardCreation(ctx)
	}
	log.Info().Msgf("[artchitect] nothing to do...")
	return nil
}

func (a *Artchitect) runCardCreation(ctx context.Context) error {
	a.state.SetState(ctx, model.StateMakingSpell)
	log.Info().Msgf("[artchitect] start card creation]")
	spell, err := a.speller.MakeSpell(ctx)
	if err != nil {
		return err
	}
	log.Info().Msgf("[artchitect] got spell: %+v", spell)
	a.state.SetState(ctx, model.StateMakingArtifact)
	card, err := a.artist.GetCard(ctx, spell)
	if err != nil {
		return err
	}
	log.Info().Msgf("[artchitect] got card: id=%d, spell_id=%d", card.ID, spell.ID)
	if err := a.notifier.NotifyNewCard(ctx, card); err != nil {
		log.Error().Err(err).Msgf("[artchitect] failed to notify new card")
	}
	a.state.SetState(ctx, model.StateMakingRest)

	return nil
}

func (a *Artchitect) runLottery(ctx context.Context, lottery model.Lottery) error {
	return a.lotteryRunner.RunLottery(ctx, lottery)
}

func (a *Artchitect) maintenance(ctx context.Context) error {
	log.Info().Msgf("[artchitect] going to maintenance")
	if err := a.lotteryRepository.InitDailyLottery(ctx); err != nil {
		return errors.Wrap(err, "[artchitect] failed InitDailyLottery from maintenance")
	}
	return nil
}
