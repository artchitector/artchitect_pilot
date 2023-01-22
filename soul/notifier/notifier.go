package notifier

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Notifier struct {
	red *redis.Client
}

func NewNotifier(red *redis.Client) *Notifier {
	return &Notifier{red: red}
}

func (n *Notifier) NotifyTick(ctx context.Context, tick int) error {
	err := n.red.Publish(ctx, model.ChannelTick, tick).Err()
	return errors.Wrap(err, "[notifier] failed to notify tick")
}

func (n *Notifier) NotifyNewCard(ctx context.Context, card model.Card) error {
	cardStr, err := json.Marshal(card)
	if err != nil {
		return errors.Wrap(err, "[notifier] failed to marshal card")
	}
	err = n.red.Publish(ctx, model.ChannelNewCard, cardStr).Err()
	return errors.Wrap(err, "[notifier] failed to notify tick")
}
