package memory

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"time"
)

var ErrNotFound = errors.New("[memory] not found")

type cache interface {
	SaveImage(ctx context.Context, cardID uint, size string, data []byte) error
	ExistsImage(ctx context.Context, ID uint, size string) (bool, error)
	GetCardImage(ctx context.Context, ID uint, size string) ([]byte, error)
}

// Memory helps get images fast (downloads from memory-server and cache locally)
type Memory struct {
	memoryURL string
	cache     cache
}

func NewMemory(memoryURL string, cache cache) *Memory {
	return &Memory{memoryURL, cache}
}

func (m *Memory) GetCardImage(ctx context.Context, cardID uint, size string) ([]byte, error) {
	if m.cache == nil {
		return nil, errors.Errorf("[memory] cache not initialized, use DownloadImage instead")
	}
	start := time.Now()

	exists, err := m.cache.ExistsImage(ctx, cardID, size)
	if err != nil {
		log.Error().Err(err).Msgf("[memory] failed check image exists %d/%s", cardID, size)
	} else if exists {
		img, err := m.cache.GetCardImage(ctx, cardID, size)
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
	if cardID == 0 {
		if dt, err := os.ReadFile(fmt.Sprintf("./files/black-%s.jpg", size)); err != nil {
			return []byte{}, errors.Wrapf(err, "[memory] failed to get black from filesystem %s", size)
		} else {
			return dt, nil
		}
	}
	// get image from remote memory server
	thousand := model.GetCardThousand(cardID)
	url := fmt.Sprintf("%s/cards/%d/card-%d-%s.jpg", m.memoryURL, thousand, cardID, size)
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "[memory] failed to get image from memory-server %s", url)
	}
	defer resp.Body.Close()
	log.Info().Msgf("STATUS %d", resp.StatusCode)
	if resp.StatusCode == http.StatusNotFound {
		return []byte{}, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return []byte{}, errors.Wrapf(err, "[memory] not OK status code(%d) from memory-server %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}

func (m *Memory) GetUnityImage(ctx context.Context, mask string, size string) ([]byte, error) {
	start := time.Now()
	defer log.Info().Msgf("[memory] get hundred image success m:%s s:%s, time: %s", mask, size, time.Now().Sub(start))
	// TODO make cached unity image
	img, err := m.downloadUnityImage(ctx, mask, size)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "[memory] failed to download hundred image m:%s s:%s", mask, size)
	}
	return img, nil
}

func (m *Memory) downloadUnityImage(ctx context.Context, mask string, size string) ([]byte, error) {
	// get image from remote memory server

	url := fmt.Sprintf("%s/unity/%s-%s.jpg", m.memoryURL, mask, size)
	log.Info().Msgf("[memory] get unity image from path %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "[memory] failed to get image from memory-server %s", url)
	}
	defer resp.Body.Close()
	log.Info().Msgf("STATUS %d", resp.StatusCode)
	if resp.StatusCode == http.StatusNotFound {
		return []byte{}, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return []byte{}, errors.Errorf("[memory] not OK status code(%d) from memory-server %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}
