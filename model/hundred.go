package model

import "gorm.io/gorm"

type Hundred struct {
	gorm.Model
	Rank    uint
	Hundred uint
}
