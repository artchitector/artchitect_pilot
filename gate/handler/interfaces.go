package handler

import (
	"context"
	"github.com/artchitector/artchitect/gate/localmodel"
	"github.com/artchitector/artchitect/model"
)

type cardsRepository interface {
	GetCard(ctx context.Context, ID uint) (model.Card, bool, error)
	GetLastCards(ctx context.Context, count uint) ([]model.Card, error)
	GetCardsRange(ctx context.Context, from uint, to uint) ([]model.Card, error)
}

type lotteryRepository interface {
	GetActiveLottery(ctx context.Context) (model.Lottery, error)
	GetLastLotteries(ctx context.Context, lastN uint) ([]model.Lottery, error)
	GetSelection(ctx context.Context) ([]uint, error)
}

type prayRepository interface {
	MakePray(ctx context.Context) (model.PrayWithQuestion, error)
	GetAnswer(ctx context.Context, prayId uint) (uint, error)
}

type cache interface {
	GetImage(ctx context.Context, ID uint, size string) ([]byte, error)
	GetCard(ctx context.Context, ID uint) (model.Card, error)
	GetLastCards(ctx context.Context, count uint) ([]model.Card, error)
}

type listener interface {
	EventChannel() (chan localmodel.Event, chan struct{})
}
