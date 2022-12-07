package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Painting struct {
	gorm.Model
	Caption string
	Bytes   sql.RawBytes
}

type PaintingPray struct {
	Caption string
}

func (p PaintingPray) Name() string {
	return EntityPainting
}

type PaintingGift struct {
	Caption  string
	Painting []byte
	Err      error
}

func (p PaintingGift) Name() string {
	return EntityPainting
}

func (p PaintingGift) Error() error {
	return p.Err
}
