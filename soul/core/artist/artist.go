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
	AddArtWatermark(originalImage image.Image, artID uint) (image.Image, error)
}

type artRepository interface {
	SaveArt(ctx context.Context, painting model.Art) (model.Art, error)
	GetLastArtPaintTime(ctx context.Context) (uint, error)
	DeleteArt(ctx context.Context, artID uint) error
}

type saver interface {
	SaveArt(ctx context.Context, artID uint, imageData []byte) error
	SaveFullsize(ctx context.Context, artID uint, imageData []byte) error
}

type Artist struct {
	engine    EngineContract
	artRepo   artRepository
	notifier  notifier
	watermark watermark
	saver     saver
}

func NewArtist(engine EngineContract, artRepository artRepository, notifier notifier, watermark watermark, saver saver) *Artist {
	return &Artist{engine, artRepository, notifier, watermark, saver}
}

func (a *Artist) GetArt(
	ctx context.Context,
	newArtID uint,
	spell model.Spell,
	artistState *model.CreationState,
) (model.Art, error) {
	log.Info().Msgf("Start get art process from artist. tags: %s, seed: %d", spell.Tags, spell.Seed)

	lastPaintingTime, err := a.artRepo.GetLastArtPaintTime(ctx)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[artist] failed to get LastPaintingTime")
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

	log.Info().Msgf("[artist] start image art with spell(id=%d)", spell.ID)
	img, err := a.engine.GetImage(ctx, spell)
	cancel()
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[artist] failed to get image-data for art")
	}

	paintTime := time.Now().Sub(paintStart)
	art := model.Art{
		Spell:     spell,
		Version:   spell.Version,
		PaintTime: uint(paintTime.Seconds()),
	}

	art.ID = newArtID

	art, err = a.artRepo.SaveArt(ctx, art)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[artist] failed to save art")
	}

	img, err = a.prepareImage(img, art.ID)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[artist] failed to prepare image")
	}

	err = a.uploadToStorage(ctx, img, art.ID)
	if err != nil {
		log.Error().Err(err).Msgf("[artist] failed to send image to storage. delete art %d", art.ID)
		if err := a.artRepo.DeleteArt(ctx, art.ID); err != nil {
			log.Error().Err(err).Msgf("[artist] failed to delete art after failed image creation (id=%d)", art.ID)
		}
		return model.Art{}, errors.Wrap(err, "[artist] failed to upload art into storage")
	}

	bts, err := a.encodeImage(img)
	if err != nil {
		log.Error().Err(err).Msgf("[artist] failed to encode image. delete art %d", art.ID)
		if err := a.artRepo.DeleteArt(ctx, art.ID); err != nil {
			log.Error().Err(err).Msgf("[artist] failed to delete art after failed image creation (id=%d)", art.ID)
		}
		return model.Art{}, errors.Wrap(err, "[artist] failed to upload art into storage")
	}

	if err := a.saver.SaveArt(ctx, art.ID, bts); err != nil {
		log.Error().Err(err).Msgf("[artist] failed to save image to saver. delete art %d", art.ID)
		if err := a.artRepo.DeleteArt(ctx, art.ID); err != nil {
			log.Error().Err(err).Msgf("[artist] failed to delete art after failed image creation (id=%d)", art.ID)
		}
		return model.Art{}, errors.Wrap(err, "[artist] failed to upload art into storage")
	}

	art.UploadedToMemory = true
	art.UploadedToStorage = true
	art, err = a.artRepo.SaveArt(ctx, art)
	if err != nil {
		// TODO need to test delete failed art without image
		if err := a.artRepo.DeleteArt(ctx, art.ID); err != nil {
			log.Error().Err(err).Msgf("[artist] failed to delete art after failed image creation (id=%d)", art.ID)
		}
		return model.Art{}, errors.Wrap(err, "[artist] failed to save art")
	}

	log.Info().Msgf("Received and saved art from artist: id=%d", art.ID)
	return art, err
}

// add watermark
func (a *Artist) prepareImage(img image.Image, artID uint) (image.Image, error) {
	var err error
	img, err = a.watermark.AddArtWatermark(img, artID)
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
func (a *Artist) uploadToStorage(ctx context.Context, img image.Image, artID uint) error {
	log.Info().Msgf("[artist] upload art %d to storage", artID)
	filename := fmt.Sprintf("art-%d.jpg", artID)
	buf := new(bytes.Buffer)

	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: model.QualityXF}); err != nil {
		return errors.Wrapf(err, "[artist] failed to encode image into jpeg with q=%d", model.QualityXF)
	}
	if err := a.saver.SaveFullsize(ctx, artID, buf.Bytes()); err != nil {
		return errors.Wrapf(err, "[artist] failed to save image to storage %s", filename)
	}
	return nil
}
