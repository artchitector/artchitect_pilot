package cache

import (
	"context"
	"github.com/artchitector/artchitect/model"
)

type artsRepository interface {
	GetArt(ctx context.Context, ID uint) (model.Art, error)
	GetLastArts(ctx context.Context, count uint) ([]model.Art, error)
}

type selectionRepository interface {
	GetSelectionLimit(ctx context.Context, limit int) ([]uint, error)
}

type memory interface {
	GetCardImage(ctx context.Context, cardID uint, size string) ([]byte, error)
}
