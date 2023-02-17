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
	GetCard(ctx context.Context, spell model.Spell, artistState *model.CreationState) (model.Card, error)
}
type speller interface {
	MakeSpell(ctx context.Context, artistState *model.CreationState) (model.Spell, error)
}
type notifier interface {
	NotifyPrehotCard(ctx context.Context, card model.Card) error
	NotifyNewCard(ctx context.Context, card model.Card) error
	NotifyCreationState(ctx context.Context, state model.CreationState) error
}

type combinator interface {
	CombineHundred(ctx context.Context, rank uint, hundred uint) error
}

// Creator used to make new Card with no input data. Used by Artchitect and Merciful
type Creator struct {
	mutex         sync.Mutex
	artist        artist
	speller       speller
	notifier      notifier
	cardTotalTime uint // in seconds
	prehotDelay   uint // in seconds
}

func NewCreator(artist artist, speller speller, notifier notifier, cardTotalTime uint, prehotDelay uint) *Creator {
	return &Creator{sync.Mutex{}, artist, speller, notifier, cardTotalTime, prehotDelay}
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
	log.Info().Msgf("[creator] unity update skipped")
	return nil
	//start := time.Now()
	//
	//if id%model.Rank100 != 0 {
	//	// update thumbs only on new 100
	//	return nil
	//}
	//
	//hundred := uint(math.Floor(float64((id-1)/model.Rank10000))) * model.Rank10000
	//log.Info().Msgf("[creator] combine rank for previous %d - %d", model.Rank10000, hundred)
	//if err := c.combinator.CombineHundred(ctx, model.Rank10000, hundred); err != nil {
	//	return errors.Wrapf(err, "[creator] failed call combinator for rank %d and card ID=%d", model.Rank10000, id)
	//}
	//
	//hundred = uint(math.Floor(float64((id-1)/model.Rank1000))) * model.Rank1000
	//log.Info().Msgf("[creator] combine rank for previous %d - %d", model.Rank1000, hundred)
	//if err := c.combinator.CombineHundred(ctx, model.Rank1000, hundred); err != nil {
	//	return errors.Wrapf(err, "[creator] failed call combinator for rank %d and card ID=%d", model.Rank1000, id)
	//}
	//
	//hundred = uint(math.Floor(float64((id-1)/model.Rank100))) * model.Rank100
	//log.Info().Msgf("[creator] combine rank for previous %d - %d", model.Rank100, hundred)
	//if err := c.combinator.CombineHundred(ctx, model.Rank100, hundred); err != nil {
	//	return errors.Wrapf(err, "[creator] failed call combinator for rank %d and card ID=%d", model.Rank100, id)
	//}
	//
	//log.Info().Msgf("[creator] update hundreds %s", time.Now().Sub(start))
	//return nil
}
