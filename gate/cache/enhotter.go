package cache

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Enhotter struct {
	artsRepository      artsRepository
	selectionRepository selectionRepository
	cache               *Cache
	memory              memory
}

func NewEnhotter(artsRepository artsRepository, selectionRepository selectionRepository, cache *Cache, memory memory) *Enhotter {
	return &Enhotter{artsRepository, selectionRepository, cache, memory}
}

func (e *Enhotter) Run(ctx context.Context) {
	go func() {
		if err := e.EnhotLastArts(ctx); err != nil {
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
			if err := e.EnhotLastArts(ctx); err != nil {
				log.Error().Err(err).Send()
			}
		case <-time.NewTicker(time.Minute * 10).C:
			if err := e.EnhotSelection(ctx); err != nil {
				log.Error().Err(err).Send()
			}
		}
	}()
}

func (e *Enhotter) EnhotLastArts(ctx context.Context) error {
	last, err := e.artsRepository.GetLastArts(ctx, 99)
	if err != nil {
		return errors.Wrap(err, "[enhotter] failed to getLastCards")
	}
	for _, art := range last {
		e.cacheArt(ctx, art)
	}
	err = e.cache.RefreshLastArts(ctx, last)
	return errors.Wrapf(err, "[enhotter] failed to RefreshLastArts in cache")
}

func (e *Enhotter) EnhotSelection(ctx context.Context) error {
	selected, err := e.selectionRepository.GetSelectionLimit(ctx, 1000)
	if err != nil {
		return errors.Wrapf(err, "[enhotter] selection get failed")
	}
	for _, artID := range selected {
		card, err := e.artsRepository.GetArt(ctx, artID)
		if err != nil {
			log.Error().Err(err).Msgf("[enhotter] failed to get card from repository id=%d", artID)
		} else {
			e.cacheArt(ctx, card)
		}
	}
	return nil
}

func (e *Enhotter) ReloadCardWithoutImage(ctx context.Context, artID uint) {
	art, err := e.artsRepository.GetArt(ctx, artID)
	if err != nil {
		log.Error().Msgf("[enhotter] failed to reload art %d", art.ID)
		return
	}
	if err := e.cache.SaveCard(ctx, art); err != nil {
		log.Error().Msgf("[enhotter] failed to saveCard %d", art.ID)
	}
	log.Info().Msgf("[enhotter] reloaded art %d", art.ID)
}

func (e *Enhotter) cacheArt(ctx context.Context, art model.Art) {
	if err := e.cache.SaveCard(ctx, art); err != nil {
		log.Error().Msgf("[enhotter] failed to saveCard %d", art.ID)
	}
	for _, size := range model.PublicSizes {
		// memory automatically cache image on get
		if _, err := e.memory.GetCardImage(ctx, art.ID, size); err != nil {
			log.Error().Err(err).Msgf("[enhotter] failed to memory.GetCardImage %d/%s", art.ID, size)
		}
	}
}
