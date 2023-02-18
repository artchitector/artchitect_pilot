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

func (sr *SpellRepository) Save(ctx context.Context, spell model.Spell) (model.Spell, error) {
	err := sr.db.Save(&spell).Error
	return spell, err
}
