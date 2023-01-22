package listener

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type cache interface {
	SaveCard(ctx context.Context, card model.Card) error
	AddLastCardID(ctx context.Context, ID uint64) error
}

type cardRepository interface {
	GetCard(ctx context.Context, ID uint) (model.Card, bool, error)
}

// Listener read incoming request from redis and do some actions
// new card saved - load it to redis
type Listener struct {
	red            *redis.Client
	cache          cache
	cardRepository cardRepository
}

func NewListener(red *redis.Client, cache cache, cardRepository cardRepository) *Listener {
	return &Listener{red, cache, cardRepository}
}

func (l *Listener) Run(ctx context.Context) error {
	subscriber := l.red.Subscribe(ctx, model.ChannelTick, model.ChannelNewCard)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("[listener] failed to receive message")
				time.Sleep(time.Second)
				continue
			}
			if err := l.handle(ctx, msg); err != nil {
				log.Error().Err(err).Msgf("[listener] failed to handle message")
			}
		}
	}
}

func (l *Listener) handle(ctx context.Context, msg *redis.Message) error {
	log.Info().Msgf("[listener] got %s event:  %s", msg.Channel, msg.Payload)
	switch msg.Channel {
	case model.ChannelTick:
		return nil
	case model.ChannelNewCard:
		return l.handleNewCard(ctx, msg)
	}
	log.Info().Msgf("%+v", msg)
	return nil
}

func (l *Listener) handleNewCard(ctx context.Context, msg *redis.Message) error {
	var card model.Card
	if err := json.Unmarshal([]byte(msg.Payload), &card); err != nil {
		return errors.Wrap(err, "[listener] failed to unmarshal new card")
	}
	log.Info().Msgf("[listener] got new card event(id=%d)", card.ID)
	card, found, err := l.cardRepository.GetCard(ctx, card.ID)
	if err != nil {
		return errors.Wrapf(err, "[listener] failed to get card id=%d", card.ID)
	} else if !found {
		return errors.Errorf("[listener] not found card id=%d", card.ID)
	}
	// card automatically saved when loaded in repository
	if err := l.cache.AddLastCardID(ctx, uint64(card.ID)); err != nil {
		return errors.Wrapf(err, "[listener] failed to append new last_card to cache with id=%d", card.ID)
	}
	return nil
}
