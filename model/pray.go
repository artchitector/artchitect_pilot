package model

import "gorm.io/gorm"

const (
	PrayStateWaiting  = "waiting"
	PrayStateRunning  = "running"
	PrayStateAnswered = "answered"
)

type Pray struct {
	gorm.Model
	Password string
	State    string
	Answer   uint
}
