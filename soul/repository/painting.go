package repository

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/model"
	"gorm.io/gorm"
)

type PaintingRepository struct {
	db *gorm.DB
}

func NewPaintingRepository(db *gorm.DB) *PaintingRepository {
	return &PaintingRepository{db}
}

func (pr *PaintingRepository) SavePainting(ctx context.Context, painting model.Painting) (model.Painting, error) {
	err := pr.db.Save(&painting).Error
	return painting, err
}
