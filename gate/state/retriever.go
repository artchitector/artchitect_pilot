package state

import (
	"context"
	"github.com/artchitector/artchitect.git/gate/model"
	"github.com/rs/zerolog"
)

type Retriever struct {
	logger zerolog.Logger
}

func NewRetriever(logger zerolog.Logger) *Retriever {
	return &Retriever{logger}
}

func (r *Retriever) CollectState(ctx context.Context) (model.State, error) {
	return model.State{
		LastPainting: model.LastPainting{ID: 5400},
	}, nil
}
