package artchitector

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

/**
Main core-service, that controls everything in system, mediator between all parts.
*/

type cloud interface {
	PrayIdea(ctx context.Context, idea model.IdeaPray) (chan model.IdeaGift, error)
}

type Artchitect struct {
	logger   zerolog.Logger
	schedule *Schedule
	cloud    cloud
}

func NewArtchitect(logger zerolog.Logger, schedule *Schedule, cld cloud) *Artchitect {
	return &Artchitect{logger, schedule, cld}
}

// Run starts main loop
func (a *Artchitect) Run(ctx context.Context) error {
	initialChan := a.schedule.MakePaintingSchedule(ctx)
	for {
		select {
		case <-initialChan:
			if err := a.paintingFlow(ctx); err != nil {
				a.logger.Error().Err(errors.Wrap(err, "failed to run painting flow")).Send()
			}
		case <-ctx.Done():
			a.logger.Info().Msg("ctx.Done caught. Artchitector shutdown")
			return nil
		}
	}
}

func (a *Artchitect) paintingFlow(ctx context.Context) error {
	idea := model.IdeaPray{}
	if ch, err := a.cloud.Wait(idea); err != nil {
		return errors.Wrap(err, "failed to wait for idea gift")
	} else {
		idea := <-ch
		a.logger.Info().Msgf("got an idea: %+v", idea)
	}
	if err := a.cloud.Pray(idea); err != nil {
		return errors.Wrap(err, "failed to pray for idea")
	}

	return nil
}
