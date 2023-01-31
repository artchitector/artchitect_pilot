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
	carder              *Carder
}

func NewEnhotter(cardRepository cardRepository, selectionRepository selectionRepository, cache *Cache, carder *Carder) *Enhotter {
	return &Enhotter{cardRepository, selectionRepository, cache, carder}
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
	for _, card := range last {
		e.carder.AddTask(card.ID, nil)
	}
	err = e.cache.RefreshLastCards(ctx, last)
	return errors.Wrapf(err, "[enhotter] failed to RefreshLastCards in cache")
}

func (e *Enhotter) EnhotSelection(ctx context.Context) error {
	selection, err := e.selectionRepository.GetSelectionLimit(ctx, 1000)
	if err != nil {
		return errors.Wrap(err, "[enhotter] failed get selection")
	}
	for _, selected := range selection[:50] {
		card, err := e.cardRepository.GetCard(ctx, selected)
		if err != nil {
			log.Error().Err(err).Msgf("[enhotter] failed to get card id=%d", selected)
		}
		if err := e.cache.SaveCard(ctx, card); err != nil {
			log.Error().Err(err).Msgf("[enhotter] failed to cache card id=%d", selected)
		}
		e.carder.AddTask(selected, nil)
	}
	for _, selected := range selection[50:200] {
		e.carder.AddTask(selected, []string{model.SizeM, model.SizeXS})
	}
	for _, selected := range selection[200:] {
		e.carder.AddTask(selected, []string{model.SizeXS})
	}
	return nil
}
