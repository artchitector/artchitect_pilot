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
	// TODO need to get image, not card. Artist is too complex
	GetArt(ctx context.Context, spell model.Spell, artistState *model.CreationState) (model.Art, error)
}
type speller interface {
	MakeSpell(ctx context.Context, artistState *model.CreationState) (model.Spell, error)
}
type notifier interface {
	NotifyPrehotCard(ctx context.Context, card model.Art) error
	NotifyNewCard(ctx context.Context, card model.Art) error
	NotifyCreationState(ctx context.Context, state model.CreationState) error
}

type unifier interface {
	UpdateUnitiesByNewCard(ctx context.Context, cardID uint) (bool, error)
}

type maxCardGetter interface {
	GetMaxArtID(ctx context.Context) (uint, error)
}

// Creator used to make new Art with no input data. Used by Artchitect and Merciful
type Creator struct {
	mutex         sync.Mutex
	artist        artist
	speller       speller
	notifier      notifier
	unifier       unifier
	maxCardGetter maxCardGetter
	cardTotalTime uint // in seconds
	prehotDelay   uint // in seconds
}

func NewCreator(
	artist artist,
	speller speller,
	notifier notifier,
	unifier unifier,
	maxCardGetter maxCardGetter,
	cardTotalTime uint,
	prehotDelay uint,
) *Creator {
	return &Creator{
		sync.Mutex{},
		artist,
		speller,
		notifier,
		unifier,
		maxCardGetter,
		cardTotalTime,
		prehotDelay,
	}
}

func (c *Creator) CreateWithoutEnjoy(ctx context.Context) (model.Art, error) {
	log.Info().Msgf("[creator] start card creation without enjoy")

	state := model.CreationState{}
	card, err := c.create(ctx, &state)

	return card, errors.Wrap(err, "[creator] failed to create card without enjoy")
}

func (c *Creator) CreateWithEnjoy(ctx context.Context) (model.Art, error) {
	log.Info().Msgf("[creator] start card creation with enjoy")
	cardStart := time.Now()

	maxCardId, err := c.maxCardGetter.GetMaxArtID(ctx)
	if err != nil {
		maxCardId = 0
	}
	state := model.CreationState{
		PreviousCardID: maxCardId,
	}

	card, err := c.create(ctx, &state)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[creator] failed to create card with enjoy")
	}

	if err := c.enjoy(ctx, &state, cardStart); err != nil {
		log.Error().Err(err).Msgf("[creator] failed enjoy :(")
	}

	return card, nil
}

func (c *Creator) create(ctx context.Context, state *model.CreationState) (model.Art, error) {
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
		return model.Art{}, err
	}
	log.Info().Msgf("[creator] got spell: %+v", spell)

	// paint card in artist
	card, err := c.artist.GetArt(ctx, spell, state)
	if err != nil {
		return model.Art{}, err
	}
	log.Info().Msgf("[creator] got card: id=%d, spell_id=%d", card.ID, spell.ID)

	// notify prehot
	if err := c.notifier.NotifyPrehotCard(ctx, card); err != nil {
		log.Error().Err(err).Msgf("[creator] failed to notify new card")
	}

	// give time to prehot cache
	<-time.After(time.Second * time.Duration(c.prehotDelay))
	// notify new card created
	if err := c.notifier.NotifyNewCard(ctx, card); err != nil {
		log.Error().Err(err).Msgf("[creator] failed to notify new card")
	}

	state.CardID = card.ID
	state.LastCardPaintTime = state.CurrentCardPaintTime
	if err := c.notifier.NotifyCreationState(ctx, *state); err != nil {
		log.Error().Err(err).Msgf("[creator] failed to notify fresh created card")
	}

	if err := c.updateUnity(ctx, card.ID); err != nil {
		log.Error().Err(err).Msgf("[creator] failed to update hundreds")
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

func (c *Creator) updateUnity(ctx context.Context, id uint) error {
	worked, err := c.unifier.UpdateUnitiesByNewCard(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "[creator] failed to updateUnities for card %d", id)
	}
	log.Info().Msgf("[creator] unifier worked=%t for card %d", worked, id)
	return nil
}
