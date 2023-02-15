package memory

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"time"
)

type cache interface {
	SaveImage(ctx context.Context, cardID uint, size string, data []byte) error
	ExistsImage(ctx context.Context, ID uint, size string) (bool, error)
	GetImage(ctx context.Context, ID uint, size string) ([]byte, error)
}

// Memory helps get images fast (downloads from memory-server and cache locally)
type Memory struct {
	memoryURL string
	cache     cache
}

func NewMemory(memoryURL string, cache cache) *Memory {
	return &Memory{memoryURL, cache}
}

func (m *Memory) GetImage(ctx context.Context, cardID uint, size string) ([]byte, error) {
	if m.cache == nil {
		return nil, errors.Errorf("[memory] cache not initialized, use DownloadImage instead")
	}
	start := time.Now()

	exists, err := m.cache.ExistsImage(ctx, cardID, size)
	if err != nil {
		log.Error().Err(err).Msgf("[memory] failed check image exists %d/%s", cardID, size)
	} else if exists {
		img, err := m.cache.GetImage(ctx, cardID, size)
		if err != nil {
			log.Error().Err(err).Msgf("[memory] failed get image fro cache %d/%s", cardID, size)
		} else {
			log.Info().Msgf("[memory] get card image success: %d/%s, cached, time:%s", cardID, size, time.Now().Sub(start))
			return img, nil
		}
	}
	img, err := m.DownloadImage(ctx, cardID, size)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "[memory] failed to download image %d/%s", cardID, size)
	}

	go func() {
		if err := m.cache.SaveImage(ctx, cardID, size, img); err != nil {
			log.Error().Err(err).Msgf("[memory] failed to cache image %d/%s", cardID, size)
		}
	}()

	log.Info().Msgf("[memory] get card image success: %d/%s, downloaded, time:%s", cardID, size, time.Now().Sub(start))
	return img, nil
}

func (m *Memory) DownloadImage(ctx context.Context, cardID uint, size string) ([]byte, error) {
	// get image from remote memory server
	thousand := model.GetCardThousand(cardID)
	url := fmt.Sprintf("%s/cards/%d/card-%d-%s.jpg", m.memoryURL, thousand, cardID, size)
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "[memory] failed to get image from memory-server %s", url)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, errors.Wrapf(err, "[memory] not OK status code(%d) from memory-server %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}