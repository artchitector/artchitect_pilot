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

type CardRepository struct {
	db      *gorm.DB
	entropy entropy
}

func NewCardRepository(db *gorm.DB, entropy entropy) *CardRepository {
	return &CardRepository{db, entropy}
}

func (pr *CardRepository) GetLastCards(ctx context.Context, count uint) ([]model.Card, error) {
	cards := make([]model.Card, 0, count)
	err := pr.db.
		Joins("Spell").
		Limit(int(count)).
		Order("cards.id desc").
		Limit(int(count)).
		Find(&cards).
		Error
	if err != nil {
		return []model.Card{}, errors.Wrapf(err, "failed to get cards count=%d", count)
	}

	return cards, err
}

func (pr *CardRepository) GetCard(ctx context.Context, ID uint) (model.Card, error) {
	card := model.Card{}
	err := pr.db.
		Joins("Spell").
		Where("cards.id = ?", ID).
		Last(&card).
		Error
	if err != nil {
		return card, errors.Wrapf(err, "[card_repository] failed to find card %d", ID)
	} else {
		return card, nil
	}
}

func (pr *CardRepository) GetCardsByRange(start uint, end uint) ([]model.Card, error) {
	var cards []model.Card
	log.Info().Msgf("[card_repo] get cards between %d and %d", start, end)
	err := pr.db.Joins("Spell").Where("cards.id between ? and ?", start, end).Find(&cards).Error
	return cards, err
}

func (pr *CardRepository) GetCards(ctx context.Context, IDs []uint) ([]model.Card, error) {
	var cards []model.Card
	err := pr.db.Joins("Spell").Where("cards.id in (?)", IDs).Find(&cards).Error
	return cards, err
}

func (pr *CardRepository) SaveCard(ctx context.Context, painting model.Card) (model.Card, error) {
	err := pr.db.Save(&painting).Error
	return painting, err
}

func (pr *CardRepository) DeleteCard(ctx context.Context, cardID uint) error {
	err := pr.db.Where("id = ?", cardID).Delete(&model.Card{}).Error
	return err
}

func (pr *CardRepository) GetCardsIDsByPeriod(ctx context.Context, start time.Time, end time.Time) ([]uint, error) {
	var ids []uint
	err := pr.db.Model(&model.Card{}).Select("id").Where("created_at between ? and ?", start, end).Find(&ids).Error
	return ids, err
}

func (pr *CardRepository) GetTotalCards(ctx context.Context) (uint, error) {
	var count uint
	err := pr.db.Select("count(id)").Model(&model.Card{}).Find(&count).Error
	return count, err
}

func (pr *CardRepository) GetOriginSelectedCard(ctx context.Context) (model.Card, error) {
	totalCards, err := pr.GetTotalCards(ctx)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[gifter] failed get total cards")
	}
	selection, err := pr.entropy.Select(ctx, totalCards)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[gifter] failed to select from origin")
	}
	card, err := pr.GetCardWithOffset(selection)
	if err != nil {
		return model.Card{}, errors.Wrapf(err, "[gifter] failed to GetCardWithOffset %d", selection-1)
	}
	return card, nil
}

func (pr *CardRepository) GetOriginSelectedCardByPeriod(ctx context.Context, start time.Time, end time.Time) (model.Card, error) {
	var total uint
	err := pr.db.Select("count(id)").Where("created_at between ? and ?", start, end).Model(&model.Card{}).Scan(&total).Error
	if err != nil {
		return model.Card{}, errors.Wrapf(err, "[card_repository] failed to get number of cards")
	}
	selection, err := pr.entropy.Select(ctx, total)
	if err != nil {
		return model.Card{}, errors.Wrapf(err, "[card_repository] failed to get selection from origin")
	}
	var card model.Card
	err = pr.db.Where("created_at between ? and ?", start, end).Limit(1).Offset(int(selection)).First(&card).Error
	if err != nil {
		return model.Card{}, errors.Wrapf(err, "[card_repository] failed to get card with offset %d", selection)
	}
	return card, nil
}

func (pr *CardRepository) GetAnyCardIDFromHundred(ctx context.Context, rank uint, start uint) (uint, error) {
	end := start + rank - 1
	log.Info().Msgf("[card_repo] GetAnyCardIDFromHundred s:%d, e:%d", start, end)
	var variants uint
	err := pr.db.Select("count(id)").Where("id between ? and ?", start, end).Model(&model.Card{}).Scan(&variants).Error
	if err != nil {
		return 0, errors.Wrapf(err, "[card_repo] failed to get variants from r:%d h:%d", rank, start)
	}
	log.Info().Msgf("[card_repo] selected max(id)=%d from r:%d h:%d", variants, rank, start)
	offset, err := pr.entropy.Select(ctx, variants)
	if err != nil {
		return 0, errors.Wrapf(err, "[card_repo] failed to get selection from origin. variants: %d", variants-start)
	}
	var targetID uint
	log.Info().Msgf("[card_repo] selecting anyCard with id between %d and %d with offset %d", start, end, offset)
	err = pr.db.
		Select("id").
		Model(&model.Card{}).
		Where("id between ? and ?", start, end).
		Order("id asc").
		Limit(1).
		Offset(int(offset)).
		Scan(&targetID).
		Error
	log.Info().Msgf("[card_repo] selected targetID %d", targetID)
	return targetID, errors.Wrapf(err, "[card_repo] failed to get targetId with start:%d,end:%d,offset:%d", start, end, offset)
}

// TODO deprecated public use, need make internal and remove usage from gifter.
// use GetOriginSelectedCard instead
func (pr *CardRepository) GetCardWithOffset(offset uint) (model.Card, error) {
	var card model.Card
	err := pr.db.
		Joins("Spell").
		Order("cards.id asc").
		Limit(1).
		Offset(int(offset)).
		Find(&card).Error
	return card, err
}

func (pr *CardRepository) GetLastCardPaintTime(ctx context.Context) (uint, error) {
	var paintTime uint
	err := pr.db.Select("paint_time").Model(&model.Card{}).Order("id desc").Limit(1).Scan(&paintTime).Error
	return paintTime, err
}

func (pr *CardRepository) GetMaxCardID(ctx context.Context) (uint, error) {
	var id uint
	err := pr.db.Select("max(id)").Model(&model.Card{}).Scan(&id).Error
	return id, err
}

func (pr *CardRepository) GetPreviousCardID(ctx context.Context, cardID uint) (uint, error) {
	var id uint
	err := pr.db.Select("id").
		Model(&model.Card{}).
		Where("id < ?", cardID).
		Order("id desc").
		Limit(1).
		Scan(&id).Error
	return id, err
}

func (pr *CardRepository) Like(ctx context.Context, cardID uint) error {
	err := pr.db.Model(model.Card{}).
		Where("id=?", cardID).
		Update("likes", gorm.Expr("likes + 1")).Error
	return err
}

func (pr *CardRepository) Unlike(ctx context.Context, cardID uint) error {
	err := pr.db.Model(model.Card{}).
		Where("id=?", cardID).
		Update("likes", gorm.Expr("case when likes > 0 then likes - 1 else 0 end")).Error
	return err
}
