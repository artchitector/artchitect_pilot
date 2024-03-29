package handler

import (
	"context"
	"github.com/artchitector/artchitect/gate/localmodel"
	"github.com/artchitector/artchitect/model"
)

type artsRepository interface {
	GetArt(ctx context.Context, ID uint) (model.Art, error)
	GetLastArts(ctx context.Context, count uint) ([]model.Art, error)
	GetArts(ctx context.Context, IDs []uint) ([]model.Art, error)
	GetArtsByRange(start uint, end uint) ([]model.Art, error)
	Like(ctx context.Context, cardID uint) error
	Unlike(ctx context.Context, cardID uint) error
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
	GetArt(ctx context.Context, ID uint) (model.Art, error)
	GetLastCards(ctx context.Context, count uint) ([]model.Art, error)
}

type listener interface {
	EventChannel() (chan localmodel.Event, chan struct{})
}

type memory interface {
	GetCardImage(ctx context.Context, cardID uint, size string) ([]byte, error)
	GetUnityImage(ctx context.Context, mask string, size string, version string) ([]byte, error)
}

type likeRepository interface {
	Like(ctx context.Context, userID uint, cardID uint) (model.Like, error)
	IsLiked(ctx context.Context, userID uint, cardID uint) (bool, error)
	GetLikes(ctx context.Context, userID uint) ([]uint, error)
}

type unityRepository interface {
	GetUnity(mask string) (model.Unity, error)
	GetRootUnities() ([]model.Unity, error)
	GetChildUnifiedUnities(parentMask string) ([]model.Unity, error)
}
type enhotter interface {
	ReloadCardWithoutImage(ctx context.Context, cardID uint)
}
