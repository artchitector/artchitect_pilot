package model

import (
	"gorm.io/gorm"
)

// State - Current system state
type State struct {
	gorm.Model
	State string // current state of artchitect
}
