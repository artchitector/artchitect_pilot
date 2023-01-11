package state

import (
	"context"
	"github.com/artchitector/artchitect/model"
	model2 "github.com/artchitector/artchitect/model"
)

type paintingRepository interface {
	GetLastCard(ctx context.Context) (model.Card, bool, error)
	GetLastCards(ctx context.Context, count uint64) ([]model.Card, error)
	GetCard(ctx context.Context, ID uint) (model.Card, bool, error)
}

type decisionRepository interface {
	GetLastDecision(ctx context.Context) (model.Decision, error)
}

type stateRepository interface {
	GetLastState(ctx context.Context) (model2.State, error)
}

type spellRepository interface {
	GetLastSpell(ctx context.Context) (model.Spell, error)
}
