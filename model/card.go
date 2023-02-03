package model

import (
	"database/sql"
	"gorm.io/gorm"
)

// TODO split card table and raw image data into separate tables and migrate database
type Card struct {
	gorm.Model
	SpellID           uint
	Spell             Spell
	Version           string // in what environment made card (tags set, version on StableDiffusion etc.)
	PaintTime         uint   // seconds, how much paint took
	UploadedToStorage bool   `gorm:"not null;default:false"` // full-size file uploaded to s3-storage
	Image             Image
}

type Image struct {
	CardID    uint         `gorm:"primaryKey"`
	Data      sql.RawBytes `json:"-"`
	Watermark bool         // is there any watermark on image?
}
