package listener

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

// Listener read incoming request from redis and do some actions
// new card saved - load it to redis
type Listener struct {
	red *redis.Client
}

func (l *Listener) Run(ctx context.Context) error {
	var red *redis.Client
	subscriber := red.Subscribe(ctx, "events")
	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("[listener] failed to receive message")
			continue
		}
		log.Info().Msgf("%+v", msg)
	}
}
