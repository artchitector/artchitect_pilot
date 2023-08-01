package gifter

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type artsRepository interface {
	GetOriginSelectedArt(ctx context.Context) (model.Art, error)
}

type artchitectBot interface {
	SendArtTo10Min(ctx context.Context, cardID uint) error
}

const (
	MaxAttempts = 10 //each photo 10 times because of context timeout
)

type Gifter struct {
	artsRepository artsRepository
	artchitectBot  artchitectBot
}

func NewGifter(
	artsRepository artsRepository,
	bot artchitectBot,
) *Gifter {
	return &Gifter{artsRepository, bot}
}

func (g *Gifter) Run(ctx context.Context) error {
	for {
		currentAttempts := 0
		for {
			currentAttempts += 1
			if currentAttempts > MaxAttempts {
				log.Info().Msgf("[gifter] max attempts (%d) exceeded", MaxAttempts)
				break
			}
			err := g.sendCard(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("[gifter] failed to send card")
			} else {
				break
			}
		}
		time.Sleep(time.Minute * 10)
	}
}

func (g *Gifter) sendCard(ctx context.Context) error {
	cardID, err := g.getArt(ctx)
	if err != nil {
		return errors.Wrap(err, "[gifter] failed to getArt")
	}

	err = g.artchitectBot.SendArtTo10Min(ctx, cardID)
	if err != nil {
		return errors.Wrapf(err, "[gifter] failed to send card id=%d to 10minchat", cardID)
	} else {
		log.Info().Msgf("[gifter] success send card id=%d to 10minchat", cardID)
		return nil
	}
}

func (g *Gifter) getArt(ctx context.Context) (uint, error) {
	art, err := g.artsRepository.GetOriginSelectedArt(ctx)
	if err != nil {
		return 0, errors.Wrapf(err, "[gifter] failed to GetOriginSelectedArt")
	}
	return art.ID, nil
}
