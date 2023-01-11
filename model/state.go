package model

import (
	"gorm.io/gorm"
	"time"
)

// State - Current system state
type State struct {
	gorm.Model
	State string // current state of artchitect
}

type LastPainting struct {
	ID      uint
	Caption string
	Spell   Spell
}

type LastDecision struct {
	ID        uint
	Result    float64
	CreatedAt time.Time
	Image     string
}

type CurrentState struct {
	CurrentState                   State
	CurrentStateDefaultTimeSeconds uint64
	LastPainting                   *LastPainting
	LastNPaintings                 []Card
	LastDecision                   *LastDecision
	LastSpell                      *Spell
}
