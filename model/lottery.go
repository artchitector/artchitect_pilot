package model

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"sort"
	"time"
)

const (
	LotteryStateWaiting  = "waiting"
	LotteryStateRunning  = "running"
	LotteryStateFinished = "finished"

	LotteryTour100 = "best-100"
	LotteryTour10  = "best-10"
	LotteryTour1   = "best-1"
)

// Lottery is process of selection best one (best ten, best 100) cards from all cards for period.
// For example, at the end of every day artchitect selects ten best works.
type Lottery struct {
	gorm.Model
	StartTime          time.Time
	CollectPeriodStart time.Time
	CollectPeriodEnd   time.Time
	State              string
	TotalWinners       uint64
	Tours              []LotteryTour
	WinnersJSON        string   `json:"-"` // as a JSON-string [123,546,232,543]
	Winners            []uint64 `gorm:"-"`
}

type LotteryTour struct {
	gorm.Model
	LotteryID   uint64
	Name        string
	MaxWinners  uint64
	WinnersJSON string // as a JSON-string [123,546,232,543]
	State       string
}

func (l *Lottery) IsFirstTour(selectedTour LotteryTour) bool {
	minTourId := selectedTour.ID
	for _, tour := range l.Tours {
		if tour.ID < minTourId {
			minTourId = tour.ID
		}
	}
	return selectedTour.ID == minTourId
}

func (l *Lottery) PreviousTour(selectedTour LotteryTour) (LotteryTour, error) {
	toursIDs := make([]uint, 0, len(l.Tours))
	for _, tour := range l.Tours {
		toursIDs = append(toursIDs, tour.ID)
	}
	sort.Slice(toursIDs, func(i, j int) bool {
		return toursIDs[i] < toursIDs[j]
	})
	var previousTourID uint
	for idx, id := range toursIDs {
		if id == selectedTour.ID {
			if idx == 0 {

				return LotteryTour{}, errors.Errorf("[lottery] failed to get PreviousTour from tour(id=%d)", selectedTour.ID)
			} else {
				previousTourID = toursIDs[idx-1]
			}
		}
	}

	for _, tour := range l.Tours {
		if tour.ID == previousTourID {
			return tour, nil
		}
	}

	return LotteryTour{}, errors.Errorf("[lottery] invalid iteration in PreviousTour from tour(id=%d)", selectedTour.ID)
}
