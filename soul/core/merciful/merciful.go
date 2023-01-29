package merciful

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type creator interface {
	CreateWithoutEnjoy(ctx context.Context) (model.Card, error)
}

type prayRepository interface {
	GetNextPray(ctx context.Context) (model.PrayWithQuestion, error)
	AnswerPray(ctx context.Context, pray model.PrayWithQuestion, answer uint) error
}

type notifier interface {
	NotifyNewCard(ctx context.Context, card model.Card) error
	NotifyCreationState(ctx context.Context, state model.CreationState) error
}

// Merciful asnwer prays
type Merciful struct {
	prayRepository prayRepository
	creator        creator
	notifier       notifier
}

func NewMerciful(prayRepository prayRepository, creator creator, notifier notifier) *Merciful {
	return &Merciful{prayRepository, creator, notifier}
}

func (m *Merciful) AnswerPray(ctx context.Context) (bool, error) {
	pray, err := m.prayRepository.GetNextPray(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // next worker will take his job
	} else if err != nil {
		return false, errors.Wrap(err, "[merciful] failed get next pray")
	}
	log.Info().Msgf("[merciful] start answering pray id=%d", pray.ID)
	card, err := m.creator.CreateWithoutEnjoy(ctx)
	if err != nil {
		return false, errors.Wrap(err, "[merciful] failed to get answer")
	}
	log.Info().Msgf("[merciful] created card id=%d for pray %d", card.ID, pray.ID)
	err = m.prayRepository.AnswerPray(ctx, pray, card.ID)
	if err != nil {
		return false, errors.Wrap(err, "[merciful] failed to save answer")
	}
	return true, nil
}
