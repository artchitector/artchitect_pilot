package repository

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/model"
	"gorm.io/gorm"
)

type SpellRepository struct {
	db *gorm.DB
}

func NewSpellRepository(db *gorm.DB) *SpellRepository {
	return &SpellRepository{db}
}

func (sr *SpellRepository) Save(ctx context.Context, spell model.Spell) (model.Spell, error) {
	err := sr.db.Save(&spell).Error
	return spell, err
}
