package cache

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Enhotter struct {
	cardRepository      cardRepository
	selectionRepository selectionRepository
	cache               *Cache
}

func NewEnhotter(cardRepository cardRepository, selectionRepository selectionRepository, cache *Cache) *Enhotter {
	return &Enhotter{cardRepository, selectionRepository, cache}
}

func (e *Enhotter) Run(ctx context.Context) {
	go func() {
		if err := e.EnhotLastCards(ctx); err != nil {
			log.Error().Err(err).Send()
		}
	}()
	go func() {
		if err := e.EnhotSelection(ctx); err != nil {
			log.Error().Err(err).Send()
		}
	}()
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.NewTicker(time.Minute * 10).C:
			if err := e.EnhotLastCards(ctx); err != nil {
				log.Error().Err(err).Send()
			}
		case <-time.NewTicker(time.Minute * 30).C:
			if err := e.EnhotSelection(ctx); err != nil {
				log.Error().Err(err).Send()
			}
		}
	}()
}

func (e *Enhotter) EnhotLastCards(ctx context.Context) error {
	last, err := e.cardRepository.GetLastCards(ctx, 99)
	if err != nil {
		return errors.Wrap(err, "[enhotter] failed to getLastCards")
	}
	err = e.cache.RefreshLastCards(ctx, last)
	return errors.Wrapf(err, "[enhotter] failed to RefreshLastCards in cache")
}

func (e *Enhotter) EnhotSelection(ctx context.Context) error {
	// TODO Придумать, как кешировать
	return nil
}
