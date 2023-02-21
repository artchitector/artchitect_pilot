package model

import (
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
	UploadedToMemory  bool   `gorm:"not null;default:false"` // file was uploaded to storage in all sizes as files
	Likes             uint   `gorm:"not null;default:0"`     // total number of likes
	Liked             bool   `gorm:"-"`                      // runtime flag, means that current user liked this image
}
