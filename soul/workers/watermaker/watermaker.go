package watermaker

import (
	"bytes"
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
)

type watermark interface {
	AddWatermark(originalImage image.Image, cardID uint) (image.Image, error)
}

// Watermaker adds watermarks to old images (temporary script)
type Watermaker struct {
	db        *gorm.DB
	watermark watermark
}

func NewWatermaker(db *gorm.DB, watermark watermark) *Watermaker {
	return &Watermaker{db: db, watermark: watermark}
}

func (w *Watermaker) Work(ctx context.Context) {
	return // already worked
	for {
		select {
		case <-ctx.Done():
			return
		default:
			//time.Sleep(time.Millisecond * 100)
			nothingToDo, err := w.makeWork()
			if err != nil {
				log.Error().Err(err).Msgf("[watermaker] found error. stop process")
				return
			} else if nothingToDo {
				log.Info().Msgf("[watermaker] no images without watermark")
				return
			}
		}
	}
}

func (w *Watermaker) makeWork() (bool, error) {
	imgData, err := w.GetNextImageWithoutWatermark()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil // nothing to do
	} else if err != nil {
		return false, errors.Wrapf(err, "[watermaker] failed to get image for card id=%d", imgData.CardID)
	}

	r := bytes.NewReader(imgData.Data)
	im, err := jpeg.Decode(r)
	if err != nil {
		return false, errors.Wrapf(err, "[watermaker] failed to decode jpeg for card %d", imgData.CardID)
	}

	im, err = w.watermark.AddWatermark(im, imgData.CardID)
	if err != nil {
		return false, errors.Wrapf(err, "[watermaker] failed to add watermark to card %d", imgData.CardID)
	}

	if err := w.SaveImage(imgData, im); err != nil {
		return false, errors.Wrapf(err, "[watermaker] failed to save image for card %d", imgData.CardID)
	}
	log.Info().Msgf("[watermaker] added watermark to image=%d", imgData.CardID)
	return false, nil
}

func (w *Watermaker) GetNextImageWithoutWatermark() (model.Image, error) {
	var img model.Image
	err := w.db.Where("watermark = false").Order("card_id asc").Limit(1).First(&img).Error
	return img, err
}

func (w *Watermaker) SaveImage(data model.Image, im image.Image) error {
	b := new(bytes.Buffer)
	if err := jpeg.Encode(b, im, &jpeg.Options{Quality: 100}); err != nil {
		return errors.Wrapf(err, "[watermaker] failed to encode jpeg. card=%d", data.CardID)
	}
	data.Data = b.Bytes()
	data.Watermark = true
	err := w.db.Save(&data).Error
	return errors.Wrapf(err, "[watermaker] failed to save image. card=%d", data.CardID)
}
