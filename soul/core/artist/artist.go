package artist

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/artchitector/artchitect/resizer"
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

type storage interface {
	Upload(ctx context.Context, filename string, file []byte) error
}

type Artist struct {
	engine    EngineContract
	cardRepo  cardRepository
	notifier  notifier
	watermark watermark
	storage   storage
}

func NewArtist(engine EngineContract, cardRepository cardRepository, notifier notifier, watermark watermark, storage storage) *Artist {
	return &Artist{engine, cardRepository, notifier, watermark, storage}
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

	img, err = a.prepareImage(img, card.ID)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[artist] failed to prepare image")
	}

	err = a.uploadToStorage(ctx, img, card.ID)
	if err != nil {
		log.Error().Err(err).Msgf("[artist] failed to send image to storage. delete card %d", card.ID)
		if err := a.cardRepo.DeleteCard(ctx, card.ID); err != nil {
			log.Error().Err(err).Msgf("[artist] failed to delete card after failed image creation (id=%d)", card.ID)
		}
		return model.Card{}, errors.Wrap(err, "[artist] failed to upload card into storage")
	}

	bts, err := a.encodeImage(img)
	if err != nil {
		log.Error().Err(err).Msgf("[artist] failed to encode image. delete card %d", card.ID)
		if err := a.cardRepo.DeleteCard(ctx, card.ID); err != nil {
			log.Error().Err(err).Msgf("[artist] failed to delete card after failed image creation (id=%d)", card.ID)
		}
		return model.Card{}, errors.Wrap(err, "[artist] failed to upload card into storage")
	}

	card.UploadedToStorage = true
	card.Image = model.Image{
		Data:      bts,
		Watermark: true,
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

// add watermark
func (a *Artist) prepareImage(img image.Image, cardID uint) (image.Image, error) {
	var err error
	img, err = a.watermark.AddWatermark(img, cardID)
	if err != nil {
		return nil, errors.Wrap(err, "[artist] failed to add watermark")
	}

	return img, nil
}

// resize image to F-size, and it will be saved to database
func (a *Artist) encodeImage(img image.Image) ([]byte, error) {
	img, err := resizer.ResizeImage(img, model.SizeF)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to resize image")
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: model.QualityF}); err != nil {
		return []byte{}, errors.Wrap(err, "[artist] failed to encode image into jpeg data")
	}
	return buf.Bytes(), nil
}

// upload image to storage with original size with quality 95
func (a *Artist) uploadToStorage(ctx context.Context, img image.Image, cardID uint) error {
	filename := fmt.Sprintf("card-%d.jpg", cardID)
	buf := new(bytes.Buffer)

	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: model.QualityXF}); err != nil {
		return errors.Wrapf(err, "[artist] failed to encode image into jpeg with q=%d", model.QualityXF)
	}
	if err := a.storage.Upload(ctx, filename, buf.Bytes()); err != nil {
		return errors.Wrapf(err, "[artist] failed to save image to storage %s", filename)
	}
	return nil
}
