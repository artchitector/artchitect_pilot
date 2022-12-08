package state

import (
	"context"
	"github.com/artchitector/artchitect.git/gate/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Retriever struct {
	logger             zerolog.Logger
	paintingRepository paintingRepository
}

func NewRetriever(logger zerolog.Logger, paintingRepository paintingRepository) *Retriever {
	return &Retriever{logger, paintingRepository}
}

func (r *Retriever) CollectState(ctx context.Context) (model.State, error) {
	lastPainting, found, err := r.paintingRepository.GetLastPainting(ctx)
	if err != nil {
		return model.State{}, errors.Wrap(err, "failed to get last painting from repo")
	}
	var lastPaintingState model.LastPainting
	if found {
		lastPaintingState = model.LastPainting{ID: lastPainting.ID, Caption: lastPainting.Caption}
	} else {
		lastPaintingState = model.LastPainting{ID: 0}
	}

	return model.State{
		LastPainting: lastPaintingState,
	}, nil
}

func (r *Retriever) GetPaintingData(ctx context.Context, paintingID uint) ([]byte, error) {
	painting, found, err := r.paintingRepository.GetPainting(ctx, paintingID)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to get painting id=%d", paintingID)
	} else if !found {
		return []byte{}, errors.Errorf("not found painting id=%d", painting)
	} else {
		return painting.Bytes, nil
	}
}
