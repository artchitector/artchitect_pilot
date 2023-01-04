package model

import (
	"gorm.io/gorm"
	"time"
)

type CurrentStateStr string

const (
	CurrentStateError      = CurrentStateStr("error")
	CurrentStateIdle       = CurrentStateStr("idle")
	CurrentStateNotWorking = CurrentStateStr("not_working")
)

// State - Current system state
type State struct {
	gorm.Model
	State string // current state of artchitect
}

type LastPainting struct {
	ID      uint
	Caption string
}

type LastDecision struct {
	Result float64
	Cdate  time.Time
	Image  string
}

type CurrentState struct {
	CurrentState CurrentStateStr
	LastPainting *LastPainting
	LastDecision *LastDecision
}
