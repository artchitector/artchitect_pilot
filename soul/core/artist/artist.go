package artist

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type cloud interface {
	Serve(ctx context.Context, prayName string) (chan model.Pray, chan model.Gift, error)
}

/*
Artist - makes paintings from prays
v 0.1 - no data, only answer after 3 seconds. For debug purposes.
*/
type Artist struct {
	logger zerolog.Logger
	cloud  cloud
}

func NewArtist(logger zerolog.Logger, cloud cloud) *Artist {
	return &Artist{logger, cloud}
}

func (a *Artist) Run(ctx context.Context) error {
	prays, gifts, err := a.cloud.Serve(ctx, model.EntityPainting)
	if err != nil {
		return errors.Wrapf(err, "failed to serve %s", model.EntityPainting)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				a.logger.Info().Msg("stop artist via ctx.Done")
				return
			case pray := <-prays:
				paintingPray, isPaintingPray := pray.(model.PaintingPray)
				if !isPaintingPray {
					a.logger.Error().Msgf("pray type is not %s", model.EntityPainting)
					continue
				}
				a.logger.Debug().Msg("artist got pray")
				data, err := a.paint(ctx)
				if err != nil {
					a.logger.Error().Err(err).Msgf("failed to paint")
					gifts <- model.PaintingGift{
						Err: errors.Wrap(err, "failed to paint"),
					}
				} else {
					a.logger.Debug().Msg("artist made gift")
					gifts <- model.PaintingGift{
						Caption:  paintingPray.Caption,
						Painting: data,
						Err:      nil,
					}
				}
			}
		}
	}()

	return nil
}

func (a *Artist) paint(ctx context.Context) ([]byte, error) {
	time.Sleep(time.Second * 3)
	content, err := os.ReadFile("files/allah.jpg")
	return content, err
}
