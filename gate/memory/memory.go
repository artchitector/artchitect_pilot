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

// Memory helps get images fast (downloads from memory-server and cache locally)
type Memory struct {
	memoryURL string
}

func NewMemory(memoryURL string) *Memory {
	return &Memory{memoryURL}
}

func (m *Memory) GetImage(ctx context.Context, cardID uint, size string) ([]byte, error) {
	start := time.Now()
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
	log.Info().Msgf("[memory] get card image success: %d/%s, time:%s", cardID, size, time.Now().Sub(start))
	return io.ReadAll(resp.Body)
}
