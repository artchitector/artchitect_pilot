package model

import "gorm.io/gorm"

type Selection struct {
	gorm.Model
	CardID    uint
	Card      Card
	LotteryID uint
	Lottery   Lottery
}
