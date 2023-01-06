package state

import (
	"context"
	model "github.com/artchitector/artchitect/model"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type State struct {
	stateRepository stateRepository
	mtx             sync.Mutex
	state           string
}

func NewState(stateRepository stateRepository) *State {
	return &State{stateRepository, sync.Mutex{}, model.StateNotWorking}
}

type stateRepository interface {
	SaveState(ctx context.Context, state model.State) (model.State, error)
}

func (s *State) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(time.Second):
			if err := s.dumpState(ctx); err != nil {
				log.Error().Err(err).Send()
			}
		}
	}
}

func (s *State) SetState(ctx context.Context, state string) {
	s.mtx.Lock()
	s.state = state
	s.mtx.Unlock()

	go func() {
		if err := s.dumpState(ctx); err != nil {
			log.Error().Err(err).Send()
		}
	}()
}

func (s *State) dumpState(ctx context.Context) error {
	_, err := s.stateRepository.SaveState(ctx, model.State{
		State: s.state,
	})
	return err
}
