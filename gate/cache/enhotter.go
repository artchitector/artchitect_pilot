package cache

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Enhotter struct {
	cardRepository      cardRepository
	selectionRepository selectionRepository
	cache               *Cache
	memory              memory
}

func NewEnhotter(cardRepository cardRepository, selectionRepository selectionRepository, cache *Cache, memory memory) *Enhotter {
	return &Enhotter{cardRepository, selectionRepository, cache, memory}
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
		case <-time.NewTicker(time.Minute * 5).C:
			if err := e.EnhotLastCards(ctx); err != nil {
				log.Error().Err(err).Send()
			}
		case <-time.NewTicker(time.Minute * 10).C:
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
	for _, card := range last {
		e.cacheCard(ctx, card)
	}
	err = e.cache.RefreshLastCards(ctx, last)
	return errors.Wrapf(err, "[enhotter] failed to RefreshLastCards in cache")
}

func (e *Enhotter) EnhotSelection(ctx context.Context) error {
	selected, err := e.selectionRepository.GetSelectionLimit(ctx, 1000)
	if err != nil {
		return errors.Wrapf(err, "[enhotter] selection get failed")
	}
	for _, cardID := range selected {
		card, err := e.cardRepository.GetCard(ctx, cardID)
		if err != nil {
			log.Error().Err(err).Msgf("[enhotter] failed to get card from repository id=%d", cardID)
		} else {
			e.cacheCard(ctx, card)
		}
	}
	return nil
}

func (e *Enhotter) ReloadCardWithoutImage(ctx context.Context, cardID uint) {
	card, err := e.cardRepository.GetCard(ctx, cardID)
	if err != nil {
		log.Error().Msgf("[enhotter] failed to reload card %d", card.ID)
		return
	}
	if err := e.cache.SaveCard(ctx, card); err != nil {
		log.Error().Msgf("[enhotter] failed to saveCard %d", card.ID)
	}
	log.Info().Msgf("[enhotter] reloaded card %d", card.ID)
}

func (e *Enhotter) cacheCard(ctx context.Context, card model.Art) {
	if err := e.cache.SaveCard(ctx, card); err != nil {
		log.Error().Msgf("[enhotter] failed to saveCard %d", card.ID)
	}
	for _, size := range model.PublicSizes {
		// memory automatically cache image on get
		if _, err := e.memory.GetCardImage(ctx, card.ID, size); err != nil {
			log.Error().Err(err).Msgf("[enhotter] failed to memory.GetCardImage %d/%s", card.ID, size)
		}
	}
}
