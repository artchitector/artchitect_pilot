package state

import (
	"context"
	"github.com/artchitector/artchitect.git/gate/model"
)

type paintingRepository interface {
	GetLastPainting(ctx context.Context) (model.Painting, bool, error)
	GetPainting(ctx context.Context, ID uint) (model.Painting, bool, error)
}

type decisionRepository interface {
	GetLastDecision(ctx context.Context) (model.Decision, error)
}

type stateRepository interface {
	GetLastState(ctx context.Context) (model.State, error)
}
