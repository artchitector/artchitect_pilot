package model

import "gorm.io/gorm"

const (
	PrayStateWaiting  = "waiting"
	PrayStateAnswered = "answered"
)

type PrayWithQuestion struct {
	gorm.Model
	State  string
	Answer uint
}
