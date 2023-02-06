package gifter

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type cardRepository interface {
	GetTotalCards(ctx context.Context) (uint, error)
	GetCardWithOffset(offset uint) (model.Card, error)
}

type origin interface {
	Select(ctx context.Context, totalVariants uint) (uint, error)
}

type artchitectBot interface {
	SendCardTo10Min(ctx context.Context, cardID uint) error
}

const (
	MaxAttempts = 10 //each photo 10 times because of context timeout
)

type Gifter struct {
	origin         origin
	cardRepository cardRepository
	artchitectBot  artchitectBot
}

func NewGifter(
	cardRepository cardRepository,
	origin origin,
	bot artchitectBot,
) *Gifter {
	return &Gifter{origin, cardRepository, bot}
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
	// TODO use cardRepository.GetOriginSelectedCard
	totalCards, err := g.cardRepository.GetTotalCards(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "[gifter] failed get total cards")
	}
	selection, err := g.origin.Select(ctx, totalCards)
	if err != nil {
		return 0, errors.Wrap(err, "[gifter] failed to select from origin")
	}
	card, err := g.cardRepository.GetCardWithOffset(selection)
	if err != nil {
		return 0, errors.Wrapf(err, "[gifter] failed to GetCardWithOffset %d", selection-1)
	}
	return card.ID, nil
}
