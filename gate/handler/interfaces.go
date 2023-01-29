package handler

import (
	"context"
	"github.com/artchitector/artchitect/gate/localmodel"
	"github.com/artchitector/artchitect/model"
)

type cardsRepository interface {
	GetCard(ctx context.Context, ID uint) (model.Card, bool, error)
	GetCardWithImage(ctx context.Context, ID uint) (model.Card, bool, error)
	GetLastCards(ctx context.Context, count uint) ([]model.Card, error)
}

type lotteryRepository interface {
	GetActiveLottery(ctx context.Context) (model.Lottery, error)
	GetLastLotteries(ctx context.Context, lastN uint) ([]model.Lottery, error)
}

type prayRepository interface {
	MakePray(ctx context.Context) (model.PrayWithQuestion, error)
	GetAnswer(ctx context.Context, prayId uint) (uint, error)
}

type selectionRepository interface {
	GetSelection(ctx context.Context) ([]uint, error)
}

type cache interface {
	GetImage(ctx context.Context, ID uint, size string) ([]byte, error)
	GetCard(ctx context.Context, ID uint) (model.Card, error)
	GetLastCards(ctx context.Context, count uint) ([]model.Card, error)
}

type listener interface {
	EventChannel() (chan localmodel.Event, chan struct{})
}
