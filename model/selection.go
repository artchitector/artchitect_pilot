package model

import "gorm.io/gorm"

type Selection struct {
	gorm.Model
	CardID    uint
	Card      Art
	LotteryID uint
	Lottery   Lottery
}
