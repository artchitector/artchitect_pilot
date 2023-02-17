package handler

import (
	"context"
	"github.com/artchitector/artchitect/gate/localmodel"
	"github.com/artchitector/artchitect/model"
)

type cardsRepository interface {
	GetCard(ctx context.Context, ID uint) (model.Card, error)
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
	GetCard(ctx context.Context, ID uint) (model.Card, error)
	GetLastCards(ctx context.Context, count uint) ([]model.Card, error)
}

type listener interface {
	EventChannel() (chan localmodel.Event, chan struct{})
}

type memory interface {
	GetCardImage(ctx context.Context, cardID uint, size string) ([]byte, error)
	GetHundredImage(ctx context.Context, rank uint, hundred uint, size string) ([]byte, error)
}

type likeRepository interface {
	Like(ctx context.Context, userID uint, cardID uint) (model.Like, error)
	IsLiked(ctx context.Context, userID uint, cardID uint) (bool, error)
	GetLikes(ctx context.Context, userID uint) ([]uint, error)
}

type hundredRepository interface {
	FindAllTenK() ([]model.Hundred, error)
	FindKList(tenKHundred uint) ([]model.Hundred, error)
	FindHList(kHundred uint) ([]model.Hundred, error)
}
