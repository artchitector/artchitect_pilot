package repository

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

// Refresher refresh cache
type Refresher struct {
	cardRepository *CardRepository
}

func NewRefresher(cardRepository *CardRepository) *Refresher {
	return &Refresher{cardRepository}
}

func (r *Refresher) StartRefreshing(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.NewTicker(time.Minute).C:
			if err := r.Refresh(ctx); err != nil {
				log.Error().Err(err).Msgf("[refresher] failed to refresh last cards")
			}
		}
	}
}

func (r *Refresher) Refresh(ctx context.Context) error {
	if _, err := r.cardRepository.GetLastCards(ctx, 100); err != nil {
		return errors.Wrapf(err, "[refresher] failed to refresh last cards")
	}
	return nil
}
