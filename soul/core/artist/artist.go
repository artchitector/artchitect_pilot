package artist

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type EngineContract interface {
	GetImage(ctx context.Context, spell model.Spell) ([]byte, error)
}

type notifier interface {
	NotifyCreationState(ctx context.Context, state model.CreationState) error
}

type cardRepository interface {
	SaveCard(ctx context.Context, painting model.Card) (model.Card, error)
	GetLastCardPaintTime(ctx context.Context) (uint, error)
}

type Artist struct {
	engine   EngineContract
	cardRepo cardRepository
	notifier notifier
}

func NewArtist(engine EngineContract, cardRepository cardRepository, notifier notifier) *Artist {
	return &Artist{engine, cardRepository, notifier}
}

func (a *Artist) GetCard(ctx context.Context, spell model.Spell, artistState *model.CreationState) (model.Card, error) {
	log.Info().Msgf("Start get painting process from artist. tags: %s, seed: %d", spell.Tags, spell.Seed)

	lastPaintingTime, err := a.cardRepo.GetLastCardPaintTime(ctx)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[artist] failed to get LastPaintingTime")
	}

	paintStart := time.Now()
	updaterCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		for {
			select {
			case <-updaterCtx.Done():
				return
			case <-time.NewTicker(time.Millisecond * 1000).C:
				artistState.LastCardPaintTime = lastPaintingTime
				artistState.CurrentCardPaintTime = uint(time.Now().Sub(paintStart).Seconds())
				if err := a.notifier.NotifyCreationState(ctx, *artistState); err != nil {
					log.Error().Err(err).Msg("[artist] failed to notify artist state")
				}
			}
		}
	}()

	log.Info().Msgf("[artist] start image painting with spell(id=%d)", spell.ID)
	data, err := a.engine.GetImage(ctx, spell)
	cancel()
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[artist] failed to get image-data for card")
	}
	paintTime := time.Now().Sub(paintStart)
	painting := model.Card{
		Spell:     spell,
		Image:     data,
		Version:   spell.Version,
		PaintTime: uint(paintTime.Seconds()),
	}

	painting, err = a.cardRepo.SaveCard(ctx, painting)
	log.Info().Msgf("Received and saved painting from artist: id=%d", painting.ID)
	return painting, err
}
