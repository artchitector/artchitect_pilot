package repository

import (
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type entropy interface {
	Select(ctx context.Context, totalVariants uint) (uint, error)
}

type ArtRepository struct {
	db      *gorm.DB
	entropy entropy
}

func NewCardRepository(db *gorm.DB, entropy entropy) *ArtRepository {
	return &ArtRepository{db, entropy}
}

func (pr *ArtRepository) GetLastArts(ctx context.Context, count uint) ([]model.Art, error) {
	arts := make([]model.Art, 0, count)
	err := pr.db.
		Joins("Spell").
		Limit(int(count)).
		Order("arts.id desc").
		Limit(int(count)).
		Find(&arts).
		Error
	if err != nil {
		return []model.Art{}, errors.Wrapf(err, "failed to get arts count=%d", count)
	}

	return arts, err
}

func (pr *ArtRepository) GetCard(ctx context.Context, ID uint) (model.Art, error) {
	art := model.Art{}
	err := pr.db.
		Joins("Spell").
		Where("arts.id = ?", ID).
		Last(&art).
		Error
	if err != nil {
		return art, errors.Wrapf(err, "[art_repository] failed to find art %d", ID)
	} else {
		return art, nil
	}
}

func (pr *ArtRepository) GetArtsByRange(start uint, end uint) ([]model.Art, error) {
	var arts []model.Art
	log.Info().Msgf("[art_repo] get arts between %d and %d", start, end)
	err := pr.db.Joins("Spell").Where("arts.id between ? and ?", start, end).Find(&arts).Error
	return arts, err
}

func (pr *ArtRepository) GetArts(ctx context.Context, IDs []uint) ([]model.Art, error) {
	var arts []model.Art
	err := pr.db.Joins("Spell").Where("arts.id in (?)", IDs).Find(&arts).Error
	return arts, err
}

func (pr *ArtRepository) SaveArt(ctx context.Context, painting model.Art) (model.Art, error) {
	err := pr.db.Save(&painting).Error
	return painting, err
}

func (pr *ArtRepository) DeleteArt(ctx context.Context, artID uint) error {
	err := pr.db.Where("id = ?", artID).Delete(&model.Art{}).Error
	return err
}

func (pr *ArtRepository) GetArtsIDsByPeriod(ctx context.Context, start time.Time, end time.Time) ([]uint, error) {
	var ids []uint
	err := pr.db.Model(&model.Art{}).Select("id").Where("created_at between ? and ?", start, end).Find(&ids).Error
	return ids, err
}

func (pr *ArtRepository) GetTotalArts(ctx context.Context) (uint, error) {
	var count uint
	err := pr.db.Select("count(id)").Model(&model.Art{}).Find(&count).Error
	return count, err
}

func (pr *ArtRepository) GetOriginSelectedArt(ctx context.Context) (model.Art, error) {
	totalArts, err := pr.GetTotalArts(ctx)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[gifter] failed get total arts")
	}
	selection, err := pr.entropy.Select(ctx, totalArts)
	if err != nil {
		return model.Art{}, errors.Wrap(err, "[gifter] failed to select from origin")
	}
	art, err := pr.GetArtWithOffset(selection)
	if err != nil {
		return model.Art{}, errors.Wrapf(err, "[gifter] failed to GetArtWithOffset %d", selection-1)
	}
	return art, nil
}

func (pr *ArtRepository) GetOriginSelectedArtByPeriod(ctx context.Context, start time.Time, end time.Time) (model.Art, error) {
	var total uint
	err := pr.db.Select("count(id)").Where("created_at between ? and ?", start, end).Model(&model.Art{}).Scan(&total).Error
	if err != nil {
		return model.Art{}, errors.Wrapf(err, "[art_repository] failed to get number of arts")
	}
	selection, err := pr.entropy.Select(ctx, total)
	if err != nil {
		return model.Art{}, errors.Wrapf(err, "[art_repository] failed to get selection from origin")
	}
	var art model.Art
	err = pr.db.Where("created_at between ? and ?", start, end).Limit(1).Offset(int(selection)).First(&art).Error
	if err != nil {
		return model.Art{}, errors.Wrapf(err, "[art_repository] failed to get art with offset %d", selection)
	}
	return art, nil
}

func (pr *ArtRepository) GetAnyCardIDFromHundred(ctx context.Context, rank uint, start uint) (uint, error) {
	end := start + rank - 1
	log.Info().Msgf("[art_repo] GetAnyCardIDFromHundred s:%d, e:%d", start, end)
	var variants uint
	err := pr.db.Select("count(id)").Where("id between ? and ?", start, end).Model(&model.Art{}).Scan(&variants).Error
	if err != nil {
		return 0, errors.Wrapf(err, "[art_repo] failed to get variants from r:%d h:%d", rank, start)
	}
	log.Info().Msgf("[art_repo] selected max(id)=%d from r:%d h:%d", variants, rank, start)
	offset, err := pr.entropy.Select(ctx, variants)
	if err != nil {
		return 0, errors.Wrapf(err, "[art_repo] failed to get selection from origin. variants: %d", variants-start)
	}
	var targetID uint
	log.Info().Msgf("[art_repo] selecting anyCard with id between %d and %d with offset %d", start, end, offset)
	err = pr.db.
		Select("id").
		Model(&model.Art{}).
		Where("id between ? and ?", start, end).
		Order("id asc").
		Limit(1).
		Offset(int(offset)).
		Scan(&targetID).
		Error
	log.Info().Msgf("[art_repo] selected targetID %d", targetID)
	return targetID, errors.Wrapf(err, "[art_repo] failed to get targetId with start:%d,end:%d,offset:%d", start, end, offset)
}

// TODO deprecated public use, need make internal and remove usage from gifter.
// use GetOriginSelectedArt instead
func (pr *ArtRepository) GetArtWithOffset(offset uint) (model.Art, error) {
	var art model.Art
	err := pr.db.
		Joins("Spell").
		Order("arts.id asc").
		Limit(1).
		Offset(int(offset)).
		Find(&art).Error
	return art, err
}

func (pr *ArtRepository) GetLastArtPaintTime(ctx context.Context) (uint, error) {
	var paintTime uint
	err := pr.db.Select("paint_time").Model(&model.Art{}).Order("id desc").Limit(1).Scan(&paintTime).Error
	return paintTime, err
}

func (pr *ArtRepository) GetMaxCardID(ctx context.Context) (uint, error) {
	var id uint
	err := pr.db.Select("case when max(id) is null then 0 else max(id) end as max_id").Model(&model.Art{}).Scan(&id).Error
	return id, err
}

func (pr *ArtRepository) GetPreviousCardID(ctx context.Context, artID uint) (uint, error) {
	var id uint
	err := pr.db.Select("id").
		Model(&model.Art{}).
		Where("id < ?", artID).
		Order("id desc").
		Limit(1).
		Scan(&id).Error
	return id, err
}

func (pr *ArtRepository) Like(ctx context.Context, artID uint) error {
	err := pr.db.Model(model.Art{}).
		Where("id=?", artID).
		Update("likes", gorm.Expr("likes + 1")).Error
	return err
}

func (pr *ArtRepository) Unlike(ctx context.Context, artID uint) error {
	err := pr.db.Model(model.Art{}).
		Where("id=?", artID).
		Update("likes", gorm.Expr("case when likes > 0 then likes - 1 else 0 end")).Error
	return err
}
