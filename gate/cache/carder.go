package cache

import (
	"context"
	"github.com/artchitector/artchitect/gate/resizer"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var defaultSizes = []string{model.SizeF, model.SizeM, model.SizeS, model.SizeXS}

type task struct {
	cardID uint
	sizes  []string
}

type Carder struct {
	cardRepository cardRepository
	cache          *Cache
	mutex          sync.Mutex
	tasks          []task
}

func NewCarder(cardRepository cardRepository, cache *Cache) *Carder {
	return &Carder{
		cardRepository,
		cache,
		sync.Mutex{},
		[]task{},
	}
}

func (c *Carder) RunWorkers(ctx context.Context, amount int) {
	for i := 0; i < amount; i++ {
		go c.runWorker(ctx, i)
	}
}

func (c *Carder) runWorker(ctx context.Context, idx int) {
	log.Info().Msgf("[carder] Running carder worker #%d", idx)
	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("[carder] stop worker #%d", idx)
			return
		default:
			t, found := c.getNextTask()
			if !found {
				time.Sleep(time.Millisecond * 100)
			} else {
				if err := c.DoneTask(ctx, t); err != nil {
					log.Error().Err(err).Msgf("[carder] failed DoneTask for card=%d", t.cardID)
				}
				log.Info().Msgf("[carder] active tasks: %d", len(c.tasks))
			}
		}
	}
}

func (c *Carder) AddTask(cardID uint, sizes []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if sizes == nil || len(sizes) == 0 {
		sizes = defaultSizes
	}

	c.tasks = append(c.tasks, task{cardID, sizes})
}

func (c *Carder) getNextTask() (task, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if len(c.tasks) == 0 {
		return task{}, false
	}
	t := c.tasks[0]
	c.tasks = c.tasks[1:]
	return t, true
}

func (c *Carder) DoneTask(ctx context.Context, t task) error {
	img, err := c.cardRepository.GetImage(ctx, t.cardID)
	if err != nil {
		return errors.Wrapf(err, "[carder] failed to get image cardID=%d", t.cardID)
	}
	for _, size := range t.sizes {
		if exists, err := c.cache.ExistsImage(ctx, t.cardID, size); err != nil {
			return errors.Wrapf(err, "[carder] failed to get existense of image cardID=%d, size=%s", t.cardID, size)
		} else if exists {
			log.Info().Msgf("[carder] image already exists cardID=%d, size=%s", t.cardID, size)
			continue
		} else {
			resized, err := resizer.Resize(img.Data, size)
			if err != nil {
				return errors.Wrapf(err, "[carder] failed to resize cardID=%d, size=%s", t.cardID, size)
			}
			if err := c.cache.SaveImage(ctx, t.cardID, size, resized); err != nil {
				return errors.Wrapf(err, "[carder] failed to save image to cache cardID=%d, size=%s", t.cardID, size)
			}
		}
	}
	return nil
}
