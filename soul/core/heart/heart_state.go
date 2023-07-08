package heart

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type notifier interface {
	NotifyHeart(ctx context.Context, state model.HeartState) error
}

type cardGiver interface {
	GetOriginSelectedCard(ctx context.Context) (model.Card, error)
}

// HeartState state updates 4-random dreams, sends them to client. Every client have same 4 rnd dreams.
type HeartState struct {
	notifier              notifier
	cardGiver             cardGiver
	rndSize               int // сколько картинок держать
	lastChangedImageIndex int

	currentState *model.HeartState
}

func NewHeartState(
	notifier notifier,
	cardGiver cardGiver,
	rndSize int,
) *HeartState {
	return &HeartState{
		notifier:              notifier,
		cardGiver:             cardGiver,
		rndSize:               rndSize,
		lastChangedImageIndex: 0,
		currentState:          nil,
	}
}

func (hs *HeartState) Run(ctx context.Context, updateTime uint) error {
	if err := hs.init(ctx); err != nil {
		return errors.Wrapf(err, "[heart_state] failed to init")
	}
	hs.notify(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("[heart_state] finished context")
			return nil
		case <-time.NewTicker(time.Second * time.Duration(updateTime)).C:
			if err := hs.work(ctx); err != nil {
				log.Error().Err(errors.Wrap(err, "[heart_state] failed to work")).Send()
			}
			hs.notify(ctx)
		}
	}
}

func (hs *HeartState) init(ctx context.Context) error {
	state := model.HeartState{Rnd: make([]uint, hs.rndSize)}
	hs.currentState = &state

	for i := 0; i < hs.rndSize; i++ {
		if err := hs.replace(ctx, i); err != nil {
			return errors.Wrapf(err, "[heart_state] failed to replace index=%d", i)
		}
	}

	return nil
}

func (hs *HeartState) work(ctx context.Context) error {
	// каждый work меняется одна из рандомных картинок, по очереди
	newIndex := hs.lastChangedImageIndex + 1
	if newIndex >= hs.rndSize {
		newIndex = 0
	}
	return hs.replace(ctx, newIndex)
}

func (hs *HeartState) replace(ctx context.Context, index int) error {
	if newCard, err := hs.cardGiver.GetOriginSelectedCard(ctx); err != nil {
		return errors.Wrapf(err, "[heart_state] failed to get new card number")
	} else {
		log.Info().Msgf("[heart_state] selected new image #%d into index %d", newCard.ID, index)
		hs.currentState.Rnd[index] = newCard.ID
		hs.lastChangedImageIndex = index
	}
	return nil
}

func (hs *HeartState) notify(ctx context.Context) {
	if err := hs.notifier.NotifyHeart(ctx, *hs.currentState); err != nil {
		log.Error().Err(err).Msgf("[heart_state] failed to notify")
	}
}
