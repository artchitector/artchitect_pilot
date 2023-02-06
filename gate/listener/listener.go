package listener

import (
	"context"
	"encoding/json"
	"github.com/artchitector/artchitect/gate/localmodel"
	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type cache interface {
	SaveCard(ctx context.Context, card model.Card) error
	PrependLastCardID(ctx context.Context, ID uint) error
}

type cardRepository interface {
	GetCard(ctx context.Context, ID uint) (model.Card, error)
}

type memory interface {
	GetImage(ctx context.Context, cardID uint, size string) ([]byte, error)
}

// Listener read incoming request from redis and do some actions
// new card saved - load it to redis
type Listener struct {
	mutex          sync.Mutex
	red            *redis.Client
	cache          cache
	cardRepository cardRepository
	memory         memory
	eventChannels  []chan localmodel.Event
}

func NewListener(red *redis.Client, cache cache, cardRepository cardRepository, memory memory) *Listener {
	return &Listener{sync.Mutex{}, red, cache, cardRepository, memory, []chan localmodel.Event{}}
}

func (l *Listener) Run(ctx context.Context) error {
	subscriber := l.red.Subscribe(
		ctx,
		model.ChannelTick,
		model.ChannelNewCard,
		model.ChannelCreation,
		model.ChannelNewSelection,
		model.ChannelLottery,
		model.ChannelPrehotCard,
	)
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
	switch msg.Channel {
	case model.ChannelTick:
	case model.ChannelCreation:
	case model.ChannelLottery:

	case model.ChannelPrehotCard:
		if err := l.handlePrehotCard(ctx, msg); err != nil {
			return errors.Wrap(err, "[listener] failed to prehot new card")
		}
		return nil // don't broadcast (it's for cache only)
	case model.ChannelNewCard:
		if err := l.handleNewCard(ctx, msg); err != nil {
			return errors.Wrap(err, "[listener] failed to handle new card")
		}
	case model.ChannelNewSelection:
		if err := l.handleNewSelection(ctx, msg); err != nil {
			return errors.Wrap(err, "[listener] failed to handle new selection")
		}

	default:
		log.Error().Msgf("[listener] unknown event %s", msg.Channel)
		return nil
	}

	err := l.broadcast(ctx, msg)
	if err != nil {
		log.Error().Err(err).Msg("[listener] failed broadcast")
	}
	return nil
}

func (l *Listener) handleNewCard(ctx context.Context, msg *redis.Message) error {
	var card model.Card
	if err := json.Unmarshal([]byte(msg.Payload), &card); err != nil {
		return errors.Wrap(err, "[listener] failed to unmarshal new card")
	}

	log.Info().Msgf("[listener] got new card event(id=%d)", card.ID)

	if err := l.cacheCard(ctx, card.ID); err != nil {
		return errors.Wrapf(err, "[listener] failed to cacheCard (id=%d)", card.ID)
	}

	if err := l.cache.PrependLastCardID(ctx, uint(card.ID)); err != nil {
		return errors.Wrapf(err, "[listener] failed to append new last_card to cache with id=%d", card.ID)
	}
	return nil
}

func (l *Listener) handlePrehotCard(ctx context.Context, msg *redis.Message) error {
	var card model.Card
	if err := json.Unmarshal([]byte(msg.Payload), &card); err != nil {
		return errors.Wrap(err, "[listener] failed to unmarshal new card")
	}
	log.Info().Msgf("[listener] prehot card(id=%d)", card.ID)
	return l.cacheCard(ctx, card.ID)
}

func (l *Listener) handleNewSelection(ctx context.Context, msg *redis.Message) error {
	var selection model.Selection
	if err := json.Unmarshal([]byte(msg.Payload), &selection); err != nil {
		return errors.Wrap(err, "[listener] failed to unmarshal new selection")
	}
	log.Info().Msgf("[listener] got new selection (id=%d)", selection.ID)
	return l.cacheCard(ctx, selection.ID)
}

func (l *Listener) cacheCard(ctx context.Context, cardID uint) error {
	card, err := l.cardRepository.GetCard(ctx, cardID)
	if err != nil {
		return errors.Wrapf(err, "[listener] failed to get card id=%d", card.ID)
	}

	for _, size := range model.PublicSizes {
		if _, err := l.memory.GetImage(ctx, cardID, size); err != nil {
			log.Error().Err(err).Msgf("[listener] failed to get image from memory %d/%s", cardID, size)
		}
	}

	if err := l.cache.SaveCard(ctx, card); err != nil {
		return errors.Wrapf(err, "[listener] failed to cache card, id=%d", card.ID)
	}

	return nil
}

func (l *Listener) EventChannel() (chan localmodel.Event, chan struct{}) {
	ch := make(chan localmodel.Event)
	l.mutex.Lock()
	defer l.mutex.Unlock()

	done := make(chan struct{})
	l.eventChannels = append(l.eventChannels, ch)
	go func(ch chan localmodel.Event) {
		<-done
		l.mutex.Lock()
		defer l.mutex.Unlock()
		var err error
		log.Info().Msgf("[listener] before chan remove: %+v", l.eventChannels)
		l.eventChannels, err = removeFromSliceByChan(l.eventChannels, ch)
		if err != nil {
			log.Error().Err(err).Msgf("[listener] failed to remove element by slice")
		}
		log.Info().Msgf("[listener] after chan remove: %+v", l.eventChannels)
	}(ch)

	return ch, done
}

func (l *Listener) broadcast(ctx context.Context, msg *redis.Message) error {
	event := localmodel.Event{
		Name:    msg.Channel,
		Payload: msg.Payload,
	}
	for _, ch := range l.eventChannels {
		l.sendEvent(ctx, ch, event)
	}
	return nil
}

func (l *Listener) sendEvent(ctx context.Context, ch chan localmodel.Event, event localmodel.Event) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("[listener] send on closed channel recovered")
			l.mutex.Lock()
			defer l.mutex.Unlock()
			var err error
			log.Info().Msgf("[listener] before %+v", l.eventChannels)
			l.eventChannels, err = removeFromSliceByChan(l.eventChannels, ch)
			if err != nil {
				log.Error().Err(err).Msgf("[listener] failed to remove closed channel from eventChannels slice")
			}
			log.Info().Msgf("[listener] after %+v", l.eventChannels)
		}
	}()
	ch <- event
}

func removeFromSliceByChan(slice []chan localmodel.Event, ch chan localmodel.Event) ([]chan localmodel.Event, error) {
	for idx, s := range slice {
		if s == ch {
			return removeFromSliceByIndex(slice, uint(idx))
		}
	}
	return nil, errors.Errorf("[listener] failed to find channel in channel list")
}

func removeFromSliceByIndex(slice []chan localmodel.Event, idx uint) ([]chan localmodel.Event, error) {
	if idx < 0 || idx >= uint(len(slice)) {
		return nil, errors.Errorf("[listener] index %d out of slice range (len=%d)", idx, len(slice))
	}
	return append(slice[:idx], slice[idx+1:]...), nil
}
