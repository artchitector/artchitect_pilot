package worker

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type saver interface {
	SaveImage(cardID uint, data []byte) error
}

type Worker struct {
	db    *gorm.DB
	saver saver
}

func NewWorker(db *gorm.DB, saver saver) *Worker {
	return &Worker{db, saver}
}

func (w *Worker) Work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("[worker] ctx.Done")
			return
		default:
			nothingToDo, err := w.copyImage(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("[worker] failed copy")
				return
			} else if nothingToDo {
				log.Error().Err(err).Msgf("[worker] nothing to do")
				return
			}
		}
	}
}

func (w *Worker) copyImage(ctx context.Context) (bool, error) {
	card, img, err := w.getNextCardAndImage()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	} else if err != nil {
		return false, err
	}
	err = w.saver.SaveImage(card.ID, img.Data)
	if err != nil {
		return false, errors.Wrapf(err, "[worker] failed to save file")
	}
	log.Info().Msgf("[worker] done with %d", card.ID)
	return false, nil
}

func (w *Worker) getNextCardAndImage() (model.Card, model.Image, error) {
	var card model.Card
	err := w.db.Where("uploaded_to_memory = false").Order("id asc").Limit(1).First(&card).Error
	if err != nil {
		return model.Card{}, model.Image{}, errors.Wrapf(err, "[worker] failed to get next card")
	}
	var image model.Image
	err = w.db.Where("card_id = ?", card.ID).First(&image).Error
	if err != nil {
		return model.Card{}, model.Image{}, errors.Wrapf(err, "[worker] failed to get image for card %d", card.ID)
	}
	return card, image, nil
}
