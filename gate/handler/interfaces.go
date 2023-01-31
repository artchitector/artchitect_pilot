package handler

import (
	"context"
	"github.com/artchitector/artchitect/gate/localmodel"
	"github.com/artchitector/artchitect/model"
)

type cardsRepository interface {
	GetCard(ctx context.Context, ID uint) (model.Card, error)
	GetImage(ctx context.Context, cardID uint) (model.Image, error)
	GetLastCards(ctx context.Context, count uint) ([]model.Card, error)
}

type lotteryRepository interface {
	GetActiveLottery(ctx context.Context) (model.Lottery, error)
	GetLastLotteries(ctx context.Context, lastN uint) ([]model.Lottery, error)
}

type prayRepository interface {
	MakePray(ctx context.Context, password string) (model.Pray, error)
	GetPrayWithPassword(ctx context.Context, prayId uint, password string) (model.Pray, error)
	GetQueueBeforePray(ctx context.Context, prayID uint) (uint, error)
}

type selectionRepository interface {
	GetSelection(ctx context.Context) ([]uint, error)
}

type cache interface {
	GetImage(ctx context.Context, ID uint, size string) ([]byte, error)
	GetCard(ctx context.Context, ID uint) (model.Card, error)
	GetLastCards(ctx context.Context, count uint) ([]model.Card, error)
	SaveImage(ctx context.Context, cardID uint, size string, data []byte) error
}

type listener interface {
	EventChannel() (chan localmodel.Event, chan struct{})
}
