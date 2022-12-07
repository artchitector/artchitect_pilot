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

type paintingRepo interface {
	SavePainting(ctx context.Context, painting model.Painting) (model.Painting, error)
}

type Artchitect struct {
	logger       zerolog.Logger
	schedule     *Schedule
	cloud        cloud
	paintingRepo paintingRepo
}

func NewArtchitect(
	logger zerolog.Logger,
	schedule *Schedule,
	cld cloud,
	paintingRepo paintingRepo,
) *Artchitect {
	return &Artchitect{
		logger,
		schedule,
		cld,
		paintingRepo,
	}
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
	gifts, err := a.cloud.Pray(ctx, model.PaintingPray{Caption: "hello world"})
	if err != nil {
		return errors.Wrap(err, "failed pray")
	}

	select {
	case <-ctx.Done():
		a.logger.Info().Msgf("stop waiting pray, ctx.Done")
		return nil
	case gft := <-gifts:
		gift, isGift := gft.(model.PaintingGift)
		if !isGift {
			return errors.Errorf("gift is not type model.PaintingGift")
		}
		if gift.Error() != nil {
			return errors.Wrap(gift.Error(), "artchitect got error instead painting")
		} else {
			a.logger.Debug().Msgf("artchitect saving gift %s", gift.Caption)
			if err := a.savePainting(ctx, gift); err != nil {
				return errors.Wrap(err, "artchitect failed to save painting")
			}
		}
		return nil
	}
}

func (a *Artchitect) savePainting(ctx context.Context, gift model.PaintingGift) error {
	painting := model.Painting{
		Caption: gift.Caption,
		Bytes:   gift.Painting,
	}
	painting, err := a.paintingRepo.SavePainting(ctx, painting)
	if err != nil {
		return errors.Wrapf(err, "failed to save painting with caption %s", painting.Caption)
	}
	a.logger.Debug().Msgf("saved painting ID=%d", painting.ID)
	return nil
}
