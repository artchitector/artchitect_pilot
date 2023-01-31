package repository

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

// Refresher refresh cache
type Refresher struct {
	cardRepository      *CardRepository
	selectionRepository *SelectionRepository
}

func NewRefresher(cardRepository *CardRepository, selectionRepository *SelectionRepository) *Refresher {
	return &Refresher{cardRepository, selectionRepository}
}

func (r *Refresher) StartRefreshing(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.NewTicker(time.Minute).C:
			if err := r.RefreshLast(ctx); err != nil {
				log.Error().Err(err).Msgf("[refresher] failed to refresh last cards")
			}
		case <-time.NewTicker(time.Minute * 10).C:
			if err := r.RefreshSelection(ctx); err != nil {
				log.Error().Err(err).Msgf("[refresher] failed to refresh selection")
			}
		}

	}
}

func (r *Refresher) RefreshLast(ctx context.Context) error {
	log.Info().Msgf("[refresher] start refresh")
	// last cards
	if _, err := r.cardRepository.GetLastCards(ctx, 100); err != nil {
		return errors.Wrapf(err, "[refresher] failed to refresh last cards")
	}
	log.Info().Msgf("[refresher] complete refresh! God bless!")
	return nil
}

func (r *Refresher) RefreshSelection(ctx context.Context) error {
	log.Info().Msgf("[refresher] god bless refresh selection - started")
	// selection
	ids, err := r.selectionRepository.GetSelection(ctx)
	if err != nil {
		return errors.Wrapf(err, "[refresher] failed to get selection")
	}
	if len(ids) > 100 {
		ids = ids[:100]
	}
	for _, id := range ids {
		if _, _, err := r.cardRepository.GetCardWithImage(ctx, uint(id)); err != nil {
			return errors.Wrapf(err, "[refresher] failed to get card id=%d", id)
		}
	}
	log.Info().Msgf("[refresher] god bless refresh selection - finished")
	return nil
}
