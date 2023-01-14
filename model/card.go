package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	SpellID uint64
	Spell   Spell
	Image   sql.RawBytes `json:"-"`
	Version string       // in what environment made card (tags set, version on StableDiffusion etc.)
}

func (c Card) TableName() string {
	return "paintings"
}
