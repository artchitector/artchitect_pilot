package merciful

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type artist interface {
	GetPainting(ctx context.Context, spell model.Spell) (model.Card, error)
}

type state interface {
	SetState(ctx context.Context, state string)
}

type speller interface {
	MakeSpell(ctx context.Context) (model.Spell, error)
}

type prayRepository interface {
	GetNextPray(ctx context.Context) (model.PrayWithQuestion, error)
	AnswerPray(ctx context.Context, pray model.PrayWithQuestion, answer uint64) error
}

// Merciful asnwer prays
type Merciful struct {
	prayRepository prayRepository
	artist         artist
	state          state
	speller        speller
}

func NewMerciful(prayRepository prayRepository, artist artist, state state, speller speller) *Merciful {
	return &Merciful{prayRepository: prayRepository, artist: artist, state: state, speller: speller}
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
	m.state.SetState(ctx, model.StateMakingSpell)
	log.Info().Msgf("[merciful] start card creation]")
	spell, err := m.speller.MakeSpell(ctx)
	if err != nil {
		return 0, err
	}
	log.Info().Msgf("[merciful] got spell: %+v", spell)
	m.state.SetState(ctx, model.StateMakingArtifact)
	card, err := m.artist.GetPainting(ctx, spell)
	if err != nil {
		return 0, err
	}
	log.Info().Msgf("[merciful] got card: id=%d, spell_id=%d", card.ID, spell.ID)
	m.state.SetState(ctx, model.StateMakingRest)

	return uint64(card.ID), nil
}
