package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex:user_card_uq" json:"-"`
	CardID uint `gorm:"uniqueIndex:user_card_uq"`
	Liked  bool
}
