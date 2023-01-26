package creator

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/rs/zerolog/log"
	"sync"
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
	mutex    sync.Mutex
	artist   artist
	speller  speller
	notifier notifier
}

func NewCreator(artist artist, speller speller, notifier notifier) *Creator {
	return &Creator{sync.Mutex{}, artist, speller, notifier}
}

func (c *Creator) Create(ctx context.Context) (model.Card, error) {
	// only one creation process at same time
	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Info().Msgf("[creator] start card creation]")

	// notify about black creation state
	state := model.CreationState{}
	if err := c.notifier.NotifyCreationState(ctx, state); err != nil {
		log.Error().Err(err).Msgf("[creator] failed notify artist state")
	}

	// generate Spell (base for card)
	spell, err := c.speller.MakeSpell(ctx, &state)
	if err != nil {
		return model.Card{}, err
	}
	log.Info().Msgf("[creator] got spell: %+v", spell)

	// paint card in artist
	card, err := c.artist.GetCard(ctx, spell, &state)
	if err != nil {
		return model.Card{}, err
	}
	log.Info().Msgf("[creator] got card: id=%d, spell_id=%d", card.ID, spell.ID)

	// notify new card created
	if err := c.notifier.NotifyNewCard(ctx, card); err != nil {
		log.Error().Err(err).Msgf("[creator] failed to notify new card")
	}

	return card, nil
}
