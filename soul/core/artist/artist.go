package artist

import (
	"bytes"
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/jpeg"
	"time"
)

type EngineContract interface {
	GetImage(ctx context.Context, spell model.Spell) (image.Image, error)
}

type notifier interface {
	NotifyCreationState(ctx context.Context, state model.CreationState) error
}

type watermark interface {
	AddWatermark(originalImage image.Image, cardID uint) (image.Image, error)
}

type cardRepository interface {
	SaveCard(ctx context.Context, painting model.Card) (model.Card, error)
	GetLastCardPaintTime(ctx context.Context) (uint, error)
	DeleteCard(ctx context.Context, cardID uint) error
}

type Artist struct {
	engine    EngineContract
	cardRepo  cardRepository
	notifier  notifier
	watermark watermark
}

func NewArtist(engine EngineContract, cardRepository cardRepository, notifier notifier, watermark watermark) *Artist {
	return &Artist{engine, cardRepository, notifier, watermark}
}

func (a *Artist) GetCard(ctx context.Context, spell model.Spell, artistState *model.CreationState) (model.Card, error) {
	log.Info().Msgf("Start get card process from artist. tags: %s, seed: %d", spell.Tags, spell.Seed)

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

	log.Info().Msgf("[artist] start image card with spell(id=%d)", spell.ID)
	img, err := a.engine.GetImage(ctx, spell)
	cancel()
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[artist] failed to get image-data for card")
	}

	paintTime := time.Now().Sub(paintStart)
	card := model.Card{
		Spell:     spell,
		Version:   spell.Version,
		PaintTime: uint(paintTime.Seconds()),
	}
	card, err = a.cardRepo.SaveCard(ctx, card)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[artist] failed to save card")
	}

	var bts []byte
	bts, err = a.prepareImage(img, card.ID)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[artist] failed to prepare image")
	}

	card.Image = model.Image{
		Data: bts,
	}
	card, err = a.cardRepo.SaveCard(ctx, card)
	if err != nil {
		// TODO need to test delete failed card without image
		if err := a.cardRepo.DeleteCard(ctx, card.ID); err != nil {
			log.Error().Err(err).Msgf("[artist] failed to delete card after failed image creation (id=%d)", card.ID)
		}
		return model.Card{}, errors.Wrap(err, "[artist] failed to save card")
	}

	log.Info().Msgf("Received and saved card from artist: id=%d", card.ID)
	return card, err
}

// decode+encode jpeg, add watermark
func (a *Artist) prepareImage(img image.Image, cardID uint) ([]byte, error) {
	var err error
	img, err = a.watermark.AddWatermark(img, cardID)
	if err != nil {
		return []byte{}, errors.Wrap(err, "[artist] failed to add watermark")
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 100}); err != nil {
		return []byte{}, errors.Wrap(err, "[artist] failed to encode image into jpeg data")
	}

	return buf.Bytes(), nil
}
