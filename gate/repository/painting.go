package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type PaintingRepository struct {
	db *gorm.DB
}

func NewPaintingRepository(db *gorm.DB) *PaintingRepository {
	return &PaintingRepository{db}
}

func (pr *PaintingRepository) GetLastPainting(ctx context.Context) (model.Painting, bool, error) {
	painting := model.Painting{}
	err := pr.db.Preload("Spell").Last(&painting).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return painting, false, nil
	} else if err != nil {
		return painting, false, errors.Wrap(err, "failed to get last painting")
	} else {
		return painting, true, nil
	}
}

func (pr *PaintingRepository) GetPainting(ctx context.Context, ID uint) (model.Painting, bool, error) {
	painting := model.Painting{}
	err := pr.db.Preload("Spell").Where("id = ?", ID).Last(&painting).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return painting, false, nil
	} else if err != nil {
		return painting, false, errors.Wrap(err, "failed to get last painting")
	} else {
		return painting, true, nil
	}
}
