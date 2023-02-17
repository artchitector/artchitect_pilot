package repository

import (
	"fmt"
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
	"strings"
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

func (ur *UnityRepository) GetRootUnities() ([]model.Unity, error) {
	var unities []model.Unity
	err := ur.db.Where("rank = ?", model.Rank10000).Order("mask desc").Find(&unities).Error
	return unities, err
}

func (ur *UnityRepository) GetChildUnities(parentMask string) ([]model.Unity, error) {
	submasks := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		submasks = append(submasks, strings.Replace(parentMask, "X", fmt.Sprintf("%d", i), 1))
	}
	var unities []model.Unity
	err := ur.db.Where("mask in (?)", submasks).Order("mask desc").Find(&unities).Error
	return unities, err
}
