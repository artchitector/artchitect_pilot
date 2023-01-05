package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Painting struct {
	gorm.Model
	SpellID uint64
	Spell   Spell
	Image   sql.RawBytes
}
