package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artchitector/artchitect/gate/resizer"
	"github.com/artchitector/artchitect/model"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

const (
	// KeyLastCards stores last cards IDs - last_cards = [1000, 999, 998, 997...]
	KeyLastCards = "last_cards"
	// KeyCard stores json data of specified card - card:1000 = {ID: 1000, ...}
	KeyCard      = "card:%d"
	KeyCardImage = "card:%d:image:%s"
)

var ErrorNotFound = errors.Errorf("[cache] not found cached data")

type Cache struct {
	mutex sync.Mutex
	rdb   *redis.Client
}

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{sync.Mutex{}, rdb}
}

func (c *Cache) Flushall(ctx context.Context) error {
	return c.rdb.FlushAll(ctx).Err()
}

func (c *Cache) GetLastCards(ctx context.Context, count uint64) ([]model.Card, error) {
	// if last card was 5-minutes old, then need to clear cache
	start := int64(0)
	stop := int64(count - 1)
	result := c.rdb.LRange(ctx, KeyLastCards, start, stop)
	if err := result.Err(); err != nil {
		return []model.Card{}, errors.Wrapf(err, "[cache] failed to get LRange %d-%d", start, stop)
	}

	ids := make([]uint64, 0, count)
	if err := result.ScanSlice(&ids); err != nil {
		return []model.Card{}, errors.Wrapf(err, "[cache] failed to scan slice")
	}

	if len(ids) < int(count) {
		return []model.Card{}, errors.Errorf("[cache] requested cards count %d, but found only %d", count, len(ids))
	}

	cards := make([]model.Card, 0, count)
	for _, id := range ids {
		if card, err := c.GetCard(ctx, id); err != nil {
			return []model.Card{}, errors.Wrapf(err, "[cache] not found cached card for last cards list. List: %+v, CardID: %d", ids, id)
		} else {
			cards = append(cards, card)
		}
	}
	log.Info().Msgf("[cache] found cards (n=%d)", len(cards))
	return cards, nil
}

func (c *Cache) GetCard(ctx context.Context, ID uint64) (model.Card, error) {
	result := c.rdb.Get(ctx, fmt.Sprintf(KeyCard, ID))
	if err := result.Err(); err == redis.Nil {
		return model.Card{}, ErrorNotFound
	} else if err != nil {
		return model.Card{}, errors.Wrapf(err, "[cache] failed to get card(id=%d)", ID)
	}
	str, err := result.Result()
	if err != nil {
		return model.Card{}, errors.Wrapf(err, "[cache] failed to get string content of card(id=%d)", ID)
	}
	var card model.Card
	if err := json.Unmarshal([]byte(str), &card); err != nil {
		return model.Card{}, errors.Wrapf(err, "[cache] failed to unmarshal content of card(id=%d)", ID)
	}
	return card, nil
}

func (c *Cache) RefreshLastCards(ctx context.Context, cards []model.Card) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ids := make([]uint64, 0, len(cards))
	for _, card := range cards {
		// each card is saved into cache
		if err := c.SaveCard(ctx, card); err != nil {
			return errors.Wrapf(err, "[cache] failed to save card(id=%d)", card.ID)
		}
		ids = append(ids, uint64(card.ID))
	}

	// delete current card list
	num, err := c.rdb.Del(ctx, KeyLastCards).Result()
	if err != nil {
		return errors.Wrapf(err, "[cache] failed to delete key %s", KeyLastCards)
	}
	log.Info().Msgf("[cache] removed %s key (n=%d)", KeyLastCards, num)

	// set new card list
	// TODO make in one operation. packed argument ids... is not working, don't understand why
	for _, id := range ids {
		err = c.rdb.RPush(ctx, KeyLastCards, id).Err()
		if err != nil {
			return errors.Wrapf(err, "[cache] failed to set last cards array into key %s", KeyLastCards)
		}
	}

	return nil
}

func (c *Cache) SaveCard(ctx context.Context, card model.Card) error {
	str, err := json.Marshal(card)
	if err != nil {
		return errors.Wrapf(err, "[cache] failed to marshal card(id=%d)", card.ID)
	}
	key := fmt.Sprintf(KeyCard, card.ID)
	err = c.rdb.Set(ctx, key, str, time.Hour).Err()
	if err != nil {
		return errors.Wrapf(err, "[cache] failed to set card into redis id=%d", card.ID)
	}

	go func() {
		// save each card size in Redis
		for _, size := range []string{model.SizeF, model.SizeM, model.SizeS, model.SizeXS} {
			exist, err := c.ExistsImage(ctx, uint64(card.ID), size)
			if err != nil {
				log.Error().Err(err).Msgf("[cache] not found existing image (id=%d, size=%s", card.ID, size)
			} else if exist {
				log.Info().Msgf("[cache] skip image resizing in cache (id=%d, size=%s", card.ID, size)
				continue
			}
			resized, err := resizer.Resize(card.Image, size)
			if err != nil {
				log.Error().Err(err).Msgf("[cache] failed to resize image into size (id=%d, size=%s)", card.ID, size)
				continue
			}
			key := fmt.Sprintf(KeyCardImage, card.ID, size)
			if err := c.rdb.Set(ctx, key, resized, time.Hour).Err(); err != nil {
				log.Error().Err(err).Msgf("[cache] failed to save size into cache (id=%d, size=%s)", card.ID, size)
				continue
			}
		}
	}()

	return nil
}

func (c *Cache) ExistsImage(ctx context.Context, ID uint64, size string) (bool, error) {
	key := fmt.Sprintf(KeyCardImage, ID, size)
	i, err := c.rdb.Exists(ctx, key).Result()
	return i > 0, err
}

func (c *Cache) GetImage(ctx context.Context, ID uint64, size string) ([]byte, error) {
	var b []byte
	key := fmt.Sprintf(KeyCardImage, ID, size)
	result := c.rdb.Get(ctx, key)
	if err := result.Err(); err != nil {
		return b, errors.Wrapf(err, "[cache] failed to get image (id=%d, size=%s)", ID, size)
	}
	return result.Bytes()
}

func (c *Cache) AddLastCardID(ctx context.Context, ID uint64) error {
	return errors.Wrapf(
		c.rdb.LPush(ctx, KeyLastCards, ID).Err(),
		"[cache] failed to append last card id=%d",
		ID,
	)
}
