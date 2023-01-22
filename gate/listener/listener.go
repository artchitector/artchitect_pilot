package listener

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Listener read incoming request from redis and do some actions
// new card saved - load it to redis
type Listener struct {
	red *redis.Client
}

func NewListener(red *redis.Client) *Listener {
	return &Listener{red: red}
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
		return errors.Wrap(err, "[listener] failed to unmarshal card")
	}
	log.Info().Msgf("[listener] got new card(id=%d)", card.ID)
	return nil
}
