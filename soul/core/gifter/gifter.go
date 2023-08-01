package gifter

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type artsRepository interface {
	GetTotalArts(ctx context.Context) (uint, error)
	GetArtWithOffset(offset uint) (model.Art, error)
}

type entropy interface {
	Select(ctx context.Context, totalVariants uint) (uint, error)
}

type artchitectBot interface {
	SendCardTo10Min(ctx context.Context, cardID uint) error
}

const (
	MaxAttempts = 10 //each photo 10 times because of context timeout
)

type Gifter struct {
	entropy        entropy
	artsRepository artsRepository
	artchitectBot  artchitectBot
}

func NewGifter(
	artsRepository artsRepository,
	entropy entropy,
	bot artchitectBot,
) *Gifter {
	return &Gifter{entropy, artsRepository, bot}
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
	cardID, err := g.getCard(ctx)
	if err != nil {
		return errors.Wrap(err, "[gifter] failed to getCard")
	}

	err = g.artchitectBot.SendCardTo10Min(ctx, cardID)
	if err != nil {
		return errors.Wrapf(err, "[gifter] failed to send card id=%d to 10minchat", cardID)
	} else {
		log.Info().Msgf("[gifter] success send card id=%d to 10minchat", cardID)
		return nil
	}
}

func (g *Gifter) getCard(ctx context.Context) (uint, error) {
	// TODO use artsRepository.GetOriginSelectedCard
	totalCards, err := g.artsRepository.GetTotalArts(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "[gifter] failed get total cards")
	}
	selection, err := g.entropy.Select(ctx, totalCards)
	if err != nil {
		return 0, errors.Wrap(err, "[gifter] failed to select from entropy")
	}
	card, err := g.artsRepository.GetArtWithOffset(selection)
	if err != nil {
		return 0, errors.Wrapf(err, "[gifter] failed to GetArtWithOffset %d", selection-1)
	}
	return card.ID, nil
}
