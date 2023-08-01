package merciful

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"sync"
)

type creator interface {
	CreateWithoutEnjoy(ctx context.Context) (model.Art, error)
}

type prayRepository interface {
	GetNextPray(ctx context.Context) (model.Pray, error)
	AnswerPray(ctx context.Context, pray model.Pray, answer uint) error
	SetPrayRunning(ctx context.Context, pray model.Pray) (model.Pray, error)
}

type notifier interface {
	NotifyNewCard(ctx context.Context, card model.Art) error
	NotifyCreationState(ctx context.Context, state model.CreationState) error
}

// Merciful asnwer prays
type Merciful struct {
	prayRepository prayRepository
	creator        creator
	notifier       notifier
	mutex          sync.Mutex
}

func NewMerciful(prayRepository prayRepository, creator creator, notifier notifier) *Merciful {
	return &Merciful{prayRepository, creator, notifier, sync.Mutex{}}
}

func (m *Merciful) AnswerPray(ctx context.Context) (bool, error) {
	m.mutex.Lock()
	m.mutex.Unlock()

	pray, err := m.prayRepository.GetNextPray(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // next worker will take his job
	} else if err != nil {
		return false, errors.Wrap(err, "[merciful] failed get next pray")
	}
	log.Info().Msgf("[merciful] start answering pray id=%d", pray.ID)
	if pray.State == model.PrayStateRunning {
		log.Info().Msgf("[merciful] rerun running pray id=%d", pray.ID)
	}
	if pray, err = m.prayRepository.SetPrayRunning(ctx, pray); err != nil {
		return false, errors.Wrapf(err, "[merciful] failed to set pray running")
	}
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
