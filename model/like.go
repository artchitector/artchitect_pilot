package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	UserID uint `json:"-"`
	CardID uint
	Liked  bool
}
