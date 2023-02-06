package notifier

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Notifier struct {
	redises map[string]*redis.Client
}

func NewNotifier(redises map[string]*redis.Client) *Notifier {
	return &Notifier{
		redises,
	}
}

func (n *Notifier) NotifyTick(ctx context.Context, tick int) error {
	err := n.publish(ctx, model.ChannelTick, tick)
	return errors.Wrap(err, "[notifier] failed to notify tick")
}

func (n *Notifier) NotifyPrehotCard(ctx context.Context, card model.Card) error {
	jsn, err := json.Marshal(card)
	if err != nil {
		return errors.Wrap(err, "[notifier] failed to marshal card")
	}
	err = n.publish(ctx, model.ChannelPrehotCard, jsn)
	return errors.Wrap(err, "[notifier] failed to notify card")
}

func (n *Notifier) NotifyNewCard(ctx context.Context, card model.Card) error {
	jsn, err := json.Marshal(card)
	if err != nil {
		return errors.Wrap(err, "[notifier] failed to marshal card")
	}
	err = n.publish(ctx, model.ChannelNewCard, jsn)
	return errors.Wrap(err, "[notifier] failed to notify card")
}

func (n *Notifier) NotifyCreationState(ctx context.Context, state model.CreationState) error {
	jsn, err := json.Marshal(state)
	if err != nil {
		return errors.Wrap(err, "[notifier] failed to marshal artist state")
	}
	err = n.publish(ctx, model.ChannelCreation, jsn)
	return errors.Wrap(err, "[notifier] failed to notify artist")
}

func (n *Notifier) NotifyNewSelection(ctx context.Context, selection model.Selection) error {
	jsn, err := json.Marshal(selection)
	if err != nil {
		return errors.Wrap(err, "[notifier] failed to marshal selection")
	}
	err = n.publish(ctx, model.ChannelNewSelection, jsn)
	return errors.Wrap(err, "[notifier] failed to notify selection")
}

func (n *Notifier) NotifyLottery(ctx context.Context, state model.LotteryState) error {
	jsn, err := json.Marshal(state)
	if err != nil {
		return errors.Wrap(err, "[notifier] failed to marshal lottery")
	}
	err = n.publish(ctx, model.ChannelLottery, jsn)
	return errors.Wrap(err, "[notifier] failed to notify lottery")
}

func (n *Notifier) publish(ctx context.Context, channel string, data interface{}) error {
	for key, r := range n.redises {
		if err := r.Publish(ctx, channel, data).Err(); err != nil {
			return errors.Wrapf(err, "[notifier] failed to publish to redis(%s)", key)
		}
	}
	return nil
}
