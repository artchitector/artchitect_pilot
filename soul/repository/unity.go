package repository

import (
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
)

type UnityRepository struct {
	db *gorm.DB
}

func NewUnityRepository(db *gorm.DB) *UnityRepository {
	return &UnityRepository{db: db}
}

func (ur *UnityRepository) GetUnity(mask string) (model.Unity, error) {
	var unity model.Unity
	err := ur.db.Where("mask = ?", mask).First(&unity).Error
	return unity, err
}
func (ur *UnityRepository) CreateUnity(mask string) (model.Unity, error) {
	unity := model.Unity{
		Mask:  mask,
		State: model.UnityStateEmpty,
	}
	err := ur.db.Save(&unity).Error
	return unity, err
}
func (ur *UnityRepository) SaveUnity(unity model.Unity) (model.Unity, error) {
	err := ur.db.Save(&unity).Error
	return unity, err
}

func (ur *UnityRepository) GetNextUnityForWork() (model.Unity, error) {
	var unity model.Unity
	err := ur.db.
		Where("state = ?", model.UnityStateEmpty).
		Or("state = ?", model.UnityStateReunification).
		Order("created_at asc").
		First(&unity).Error
	return unity, err
}
