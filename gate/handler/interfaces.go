package handler

import (
	"context"
	"github.com/artchitector/artchitect.git/gate/model"
)

type retriever interface {
	CollectState(ctx context.Context) (model.State, error)
	GetPaintingData(ctx context.Context, paintingID uint) ([]byte, error)
}
