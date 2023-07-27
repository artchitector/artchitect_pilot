package streamer

import (
	"context"
	"encoding/base64"
	"github.com/artchitector/artchitect/model"
	"github.com/rs/zerolog/log"
	"sync"
)

type Streamer struct {
	mutex    sync.Mutex
	channels map[string][]chan []byte
}

func NewStreamer() *Streamer {
	return &Streamer{mutex: sync.Mutex{}, channels: make(map[string][]chan []byte)}
}

func (s *Streamer) OnPhaseEvent(ctx context.Context, state model.EntropyState) {
	data, err := base64.StdEncoding.DecodeString(state.Image)
	if err != nil {
		log.Error().Msgf("[streamer] failed base64 decode")
	}

	log.Info().Msgf("[streamer] got message. state=%s, imgSize=%d", state.Phase, len(data))

	s.mutex.Lock()
	chs, ok := s.channels[state.Phase]
	s.mutex.Unlock()

	if !ok {
		s.channels[state.Phase] = make([]chan []byte, 0)
		return
	}

	// Нужно получить изображение из редиса и отправить его по всем каналам, ожидающим его
	for _, ch := range chs {
		ch <- data
	}
}

func (s *Streamer) GetStream(ctx context.Context, phase string) (chan []byte, func()) {
	ch := s.makeChan(ctx, phase)
	chCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-chCtx.Done()
		log.Info().Msgf("[streamer] Channel %+v context Done. Remove chan", ch)
		s.mutex.Lock()
		defer s.mutex.Unlock()
		// need delete old channel from queue
		oldChannels := s.channels[phase]
		for i, oldCh := range oldChannels {
			if oldCh == ch {
				s.channels[phase] = append(oldChannels[:i], oldChannels[i+1:]...)
				close(oldCh)
				break
			}
		}
		log.Info().Msgf("[streamer] New channels: %+v", s.channels)
	}()
	return ch, cancel
}

func (s *Streamer) makeChan(ctx context.Context, phase string) chan []byte {
	ch := make(chan []byte)
	s.channels[phase] = append(s.channels[phase], ch)
	log.Info().Msgf("[streamer] started new channel %+v:%s. Current channels: %+v", ch, phase, s.channels)
	return ch
}
