package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
)

type PrayRepository struct {
	db *gorm.DB
}

func NewPrayRepository(db *gorm.DB) *PrayRepository {
	return &PrayRepository{db: db}
}

func (pr *PrayRepository) MakePray(ctx context.Context, password string) (model.Pray, error) {
	passEncrypted := encrypt(password)
	pray := model.Pray{
		Password: passEncrypted,
		State:    model.PrayStateWaiting,
		Answer:   0,
	}
	err := pr.db.Create(&pray).Error
	return pray, err
}

func (pr *PrayRepository) GetPrayWithPassword(ctx context.Context, prayID uint, password string) (model.Pray, error) {
	passEncrypted := encrypt(password)
	var pray model.Pray
	err := pr.db.Model(&model.Pray{}).
		Where("id = ?", prayID).
		Where("password = ?", passEncrypted).
		First(&pray).
		Error
	return pray, err
}

func (pr *PrayRepository) GetQueueBeforePray(ctx context.Context, prayID uint) (uint, error) {
	var queue uint
	err := pr.db.Model(&model.Pray{}).
		Select("count(id)").
		Where("state = ? or state = ?", model.PrayStateWaiting, model.PrayStateRunning).
		Where("id < ?", prayID).
		Scan(&queue).Error
	return queue, err
}

func encrypt(pass string) string {
	hash := md5.Sum([]byte(pass))
	return hex.EncodeToString(hash[:])
}
