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
	Pray(ctx context.Context, pray model.Pray) (chan model.Gift, error)
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
	scheduleChan := a.schedule.MakePaintingSchedule(ctx)
	for {
		select {
		case <-scheduleChan:
			a.logger.Debug().Msg("starting painting flow")
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
	gifts, err := a.cloud.Pray(ctx, model.Pray{Name: model.EntityPainting})
	if err != nil {
		return errors.Wrap(err, "failed pray")
	}

	select {
	case <-ctx.Done():
		a.logger.Info().Msgf("stop waiting pray, ctx.Done")
		return nil
	case gift := <-gifts:
		if gift.Error != nil {
			return errors.Wrap(gift.Error, "artchitect got error instead painting")
		} else {
			a.logger.Debug().Msgf("artchitect got gift %+v", gift)
		}
		return nil
	}
}
