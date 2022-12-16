package artchitector

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

type Schedule struct {
	logger zerolog.Logger
}

func NewSchedule(logger zerolog.Logger) *Schedule {
	return &Schedule{logger}
}

func (s *Schedule) MakePaintingSchedule(ctx context.Context) chan struct{} {
	ch := make(chan struct{})
	tick := time.NewTicker(time.Second * 60)
	go func() {
		for {
			select {
			case <-tick.C:
				// request new painting via chan, that used by artchitector
				ch <- struct{}{}
			case <-ctx.Done():
				s.logger.Info().Msg("ctx.Done caught")
				return
			}
		}
	}()
	return ch
}
