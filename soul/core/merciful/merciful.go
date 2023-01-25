package merciful

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type artist interface {
	GetCard(ctx context.Context, spell model.Spell, artistState *model.ArtistState) (model.Card, error)
}

type state interface {
	SetState(ctx context.Context, state string)
}

type speller interface {
	MakeSpell(ctx context.Context, artistState *model.ArtistState) (model.Spell, error)
}

type prayRepository interface {
	GetNextPray(ctx context.Context) (model.PrayWithQuestion, error)
	AnswerPray(ctx context.Context, pray model.PrayWithQuestion, answer uint64) error
}

type notifier interface {
	NotifyNewCard(ctx context.Context, card model.Card) error
	NotifyArtistState(ctx context.Context, state model.ArtistState) error
}

// Merciful asnwer prays
type Merciful struct {
	prayRepository prayRepository
	artist         artist
	state          state
	speller        speller
	notifier       notifier
}

func NewMerciful(prayRepository prayRepository, artist artist, state state, speller speller, notifier notifier) *Merciful {
	return &Merciful{prayRepository, artist, state, speller, notifier}
}

func (m *Merciful) AnswerPray(ctx context.Context) (bool, error) {
	pray, err := m.prayRepository.GetNextPray(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // next worker will take his job
	} else if err != nil {
		return false, errors.Wrap(err, "[merciful] failed get next pray")
	}
	// need to make a picture and get its ID
	cardId, err := m.getAnswer(ctx)
	if err != nil {
		return false, errors.Wrap(err, "[merciful] failed to get answer")
	}
	err = m.prayRepository.AnswerPray(ctx, pray, cardId)
	if err != nil {
		return false, errors.Wrap(err, "[merciful] failed to save answer")
	}
	return true, nil
}

func (m *Merciful) getAnswer(ctx context.Context) (uint64, error) {
	log.Info().Msgf("[merciful] start card creation]")
	state := model.ArtistState{}
	if err := m.notifier.NotifyArtistState(ctx, state); err != nil {
		log.Error().Err(err).Msgf("[merciful] failed notify artist state")
	}
	spell, err := m.speller.MakeSpell(ctx, &state)
	if err != nil {
		return 0, err
	}
	log.Info().Msgf("[merciful] got spell: %+v", spell)
	card, err := m.artist.GetCard(ctx, spell, &state)
	if err != nil {
		return 0, err
	}
	log.Info().Msgf("[merciful] got card: id=%d, spell_id=%d", card.ID, spell.ID)
	if err := m.notifier.NotifyNewCard(ctx, card); err != nil {
		log.Error().Err(err).Msgf("[merciful] failed to notify new card")
	}

	return uint64(card.ID), nil
}
