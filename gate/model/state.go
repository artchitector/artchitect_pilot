package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	StateError          = "error"
	StateNotWorking     = "not_working"
	StateMakingSpell    = "making_spell"
	StateMakingArtifact = "making_artifact"
	StateMakingRest     = "making_rest"
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
	LastDecision                   *LastDecision
	LastSpell                      *Spell
}
