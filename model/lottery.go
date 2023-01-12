package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	LotteryStateWaiting  = "waiting"
	LotteryStateRunning  = "running"
	LotteryStateFinished = "finished"

	LotteryTypeDaily = "daily"
)

// Lottery is process of selection best one (best ten, best 100) cards from all cards for period.
// For example, at the end of every day artchitect selects ten best works.
type Lottery struct {
	gorm.Model
	Name               string
	Type               string
	StartTime          time.Time
	CollectPeriodStart time.Time
	CollectPeriodEnd   time.Time
	Started            time.Time
	Finished           time.Time
	State              string
	TotalWinners       uint64
	WinnersJSON        string   `json:"-"` // as a JSON-string [123,546,232,543]
	Winners            []uint64 `gorm:"-"`
}
