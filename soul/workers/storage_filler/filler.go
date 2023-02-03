package storage_filler

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/artchitector/artchitect/resizer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type storage interface {
	Upload(ctx context.Context, filename string, file []byte) error
}

type Filler struct {
	db      *gorm.DB
	storage storage
}

func NewFiller(db *gorm.DB, storage storage) *Filler {
	return &Filler{db: db, storage: storage}
}

func (f *Filler) Work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			nothingToDo, err := f.makeWork(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("[storage_filler] found error. stop process")
				return
			} else if nothingToDo {
				log.Info().Msgf("[storage_filler] no images not uploaded. stop worker")
				return
			}
			return
		}
	}
}

func (f *Filler) makeWork(ctx context.Context) (bool, error) {
	card, image, err := f.GetNextNotUploadedCard()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	} else if err != nil {
		return false, errors.Wrap(err, "[storage_filler] failed to get card and image")
	}

	// Выложить оригинал изображения в s3
	filename := fmt.Sprintf("card-%d.jpg", card.ID)
	if err := f.storage.Upload(ctx, filename, image.Data); err != nil {
		return false, errors.Wrapf(err, "[storage_filler] failed to upload image to s3 for card %d", card.ID)
	}

	if card.ID >= 571 {
		// Пережать изображение в размер f (если ID>=571, до этого карточки были 512х512)
		resized, err := resizer.ResizeBytes(image.Data, model.SizeF)
		if err != nil {
			return false, errors.Wrapf(err, "[storage_filler] failed to resize image, card %d", card.ID)
		}
		image.Data = resized

		// Сохранить новые биты в Image
		_, err = f.SaveImage(image)
		if err != nil {
			return false, errors.Wrapf(err, "[storage_filler] failed to save image, card %d", card.ID)
		}
	} else {
		log.Info().Msgf("[storage_filler} skip image (card_id=%d) resize, very old", card.ID)
	}

	// Пометить карточку как uploaded_to_storage=true
	card.UploadedToStorage = true
	card, err = f.SaveCard(card)
	if err != nil {
		return false, errors.Wrapf(err, "[storage_filler] failed to save card, card %d", card.ID)
	}
	log.Info().Msgf("[storage_filler] successful uploaded and resized card=%d", card.ID)
	return false, nil
}

func (f *Filler) GetNextNotUploadedCard() (model.Card, model.Image, error) {
	var card model.Card
	var image model.Image
	err := f.db.Where("uploaded_to_storage = false").Order("id asc").Limit(1).First(&card).Error
	if err != nil {
		return model.Card{}, model.Image{}, err
	}
	err = f.db.Where("card_id = ?", card.ID).First(&image).Error
	if err != nil {
		return model.Card{}, model.Image{}, err
	}
	return card, image, nil
}

func (f *Filler) SaveImage(img model.Image) (model.Image, error) {
	err := f.db.Save(&img).Error
	return img, err
}

func (f *Filler) SaveCard(card model.Card) (model.Card, error) {
	err := f.db.Save(&card).Error
	return card, err
}
