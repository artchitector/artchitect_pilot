package creator

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type artist interface {
	GetCard(ctx context.Context, spell model.Spell, artistState *model.CreationState) (model.Card, error)
}
type speller interface {
	MakeSpell(ctx context.Context, artistState *model.CreationState) (model.Spell, error)
}
type notifier interface {
	NotifyNewCard(ctx context.Context, card model.Card) error
	NotifyCreationState(ctx context.Context, state model.CreationState) error
}

// Creator used to make new Card with no input data. Used by Artchitect and Merciful
type Creator struct {
	mutex         sync.Mutex
	artist        artist
	speller       speller
	notifier      notifier
	cardTotalTime uint // in seconds
}

func NewCreator(artist artist, speller speller, notifier notifier, cardTotalTime uint) *Creator {
	return &Creator{sync.Mutex{}, artist, speller, notifier, cardTotalTime}
}

func (c *Creator) CreateWithoutEnjoy(ctx context.Context) (model.Card, error) {
	log.Info().Msgf("[creator] start card creation without enjoy")

	state := model.CreationState{}
	card, err := c.create(ctx, &state)

	return card, errors.Wrap(err, "[creator] failed to create card without enjoy")
}

func (c *Creator) CreateWithEnjoy(ctx context.Context) (model.Card, error) {
	log.Info().Msgf("[creator] start card creation with enjoy")
	cardStart := time.Now()

	state := model.CreationState{}

	card, err := c.create(ctx, &state)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[creator] failed to create card with enjoy")
	}

	if err := c.enjoy(ctx, &state, cardStart); err != nil {
		log.Error().Err(err).Msgf("[creator] failed enjoy :(")
	}

	return card, nil
}

func (c *Creator) create(ctx context.Context, state *model.CreationState) (model.Card, error) {
	// only one creation process at same time
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// notify about black creation state
	if err := c.notifier.NotifyCreationState(ctx, *state); err != nil {
		log.Error().Err(err).Msgf("[creator] failed notify artist state")
	}

	// generate Spell (base for card)
	spell, err := c.speller.MakeSpell(ctx, state)
	if err != nil {
		return model.Card{}, err
	}
	log.Info().Msgf("[creator] got spell: %+v", spell)

	// paint card in artist
	card, err := c.artist.GetCard(ctx, spell, state)
	if err != nil {
		return model.Card{}, err
	}
	log.Info().Msgf("[creator] got card: id=%d, spell_id=%d", card.ID, spell.ID)

	// notify new card created
	if err := c.notifier.NotifyNewCard(ctx, card); err != nil {
		log.Error().Err(err).Msgf("[creator] failed to notify new card")
	}

	state.CardID = card.ID
	state.LastCardPaintTime = state.CurrentCardPaintTime
	if err := c.notifier.NotifyCreationState(ctx, *state); err != nil {
		log.Error().Err(err).Msgf("[creator] failed to notify fresh created card")
	}

	return card, nil
}

// wait till 48 seconds, because every card creates minimum 48 seconds
func (c *Creator) enjoy(ctx context.Context, state *model.CreationState, cardStart time.Time) error {
	enjoyStart := time.Now()
	cardEnd := cardStart.Add(time.Second * time.Duration(c.cardTotalTime))
	if enjoyStart.After(cardEnd) {
		log.Warn().Msgf("[creator] card was too slow, no enjoy!")
		return nil // card is too slow
	}
	secondsLeft := cardEnd.Sub(enjoyStart).Seconds()
	log.Info().Msgf("[creator] enjoy for %f seconds", secondsLeft)

	state.EnjoyTime = uint(secondsLeft)
	updaterCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		for {
			select {
			case <-updaterCtx.Done():
				return
			case <-time.NewTicker(time.Second).C:
				state.CurrentEnjoyTime = uint(time.Now().Sub(enjoyStart).Seconds())
				if err := c.notifier.NotifyCreationState(ctx, *state); err != nil {
					log.Error().Err(err).Msgf("[creator] failed to notify enjoy time")
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-time.After(time.Duration(secondsLeft) * time.Second):
		return nil // wait
	}
}