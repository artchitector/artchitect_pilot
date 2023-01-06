package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
)

type SpellRepository struct {
	db *gorm.DB
}

func NewSpellRepository(db *gorm.DB) *SpellRepository {
	return &SpellRepository{db}
}

func (dr *SpellRepository) GetLastSpell(ctx context.Context) (model.Spell, error) {
	spell := model.Spell{}
	err := dr.db.Last(&spell).Error
	return spell, err
}
