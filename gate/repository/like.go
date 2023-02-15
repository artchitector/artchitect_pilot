package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db}
}

func (lr *LikeRepository) Like(ctx context.Context, userID uint, cardID uint) (model.Like, error) {
	var like model.Like
	err := lr.db.Where("card_id = ?", cardID).Where("user_id = ?", userID).Limit(1).First(&like).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Like{}, errors.Wrapf(err, "[like_repo] failed to get like %d-%d", userID, cardID)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		like = model.Like{
			UserID: userID,
			CardID: cardID,
			Liked:  false,
		}
	}
	like.Liked = !like.Liked
	err = lr.db.Save(&like).Error
	return like, err
}

func (lr *LikeRepository) IsLiked(ctx context.Context, userID uint, cardID uint) (bool, error) {
	var like model.Like
	err := lr.db.Where("card_id = ?", cardID).
		Where("user_id = ?", userID).
		Limit(1).
		First(&like).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return like.Liked, nil
	}
}

func (lr *LikeRepository) GetLikes(ctx context.Context, userID uint) ([]uint, error) {
	var ids []uint
	err := lr.db.Select("card_id").
		Model(&model.Like{}).
		Where("user_id = ?", userID).
		Where("liked = true").
		Order("created_at DESC").
		Scan(&ids).Error
	return ids, err
}
