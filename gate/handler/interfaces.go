package handler

import (
	"context"
	"github.com/artchitector/artchitect/model"
)

type retriever interface {
	CollectState(ctx context.Context) (model.CurrentState, error)
	GetPaintingData(ctx context.Context, paintingID uint) ([]byte, error)
}

type paintingsRepository interface {
	GetLastPaintings(ctx context.Context, count uint64) ([]model.Painting, error)
}
