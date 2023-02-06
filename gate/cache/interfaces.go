package cache

import (
	"context"
	"github.com/artchitector/artchitect/model"
)

type cardRepository interface {
	GetCard(ctx context.Context, ID uint) (model.Card, error)
	GetLastCards(ctx context.Context, count uint) ([]model.Card, error)
}

type selectionRepository interface {
	GetSelectionLimit(ctx context.Context, limit int) ([]uint, error)
}
