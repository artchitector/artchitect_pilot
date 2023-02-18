package repository

import (
	"fmt"
	"github.com/artchitector/artchitect/model"
	"gorm.io/gorm"
	"log"
	"math"
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

func (ur *UnityRepository) GetChildUnifiedUnities(parentMask string) ([]model.Unity, error) {
	submasks := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		submasks = append(submasks, strings.Replace(parentMask, "X", fmt.Sprintf("%d", i), 1))
	}
	var unities []model.Unity
	err := ur.db.Where("mask in (?)", submasks).Where("state = ?", model.UnityStateUnified).Order("mask desc").Find(&unities).Error
	return unities, err
}

func (ur *UnityRepository) CreateUnity(mask string) (model.Unity, error) {
	cnt := strings.Count(mask, "X")
	var rank uint
	switch cnt {
	case 4:
		rank = model.Rank10000
	case 3:
		rank = model.Rank1000
	case 2:
		rank = model.Rank100
	}

	unity := model.Unity{
		Mask:  mask,
		Rank:  rank,
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

func (ur *UnityRepository) CreateUnityByCard(cardID uint, rank uint) (model.Unity, error) {
	mask := calculateMask(cardID, rank)
	un := model.Unity{
		Mask:  mask,
		Rank:  rank,
		State: model.UnityStateEmpty,
	}
	err := ur.db.Save(&un).Error
	return un, err
}

func (ur *UnityRepository) GetUnityByCard(cardID uint, rank uint) (model.Unity, error) {
	mask := calculateMask(cardID, rank)
	var un model.Unity
	err := ur.db.Where("mask = ?", mask).First(&un).Error
	return un, err
}

func calculateMask(cardID uint, rank uint) string {
	normalized := int(math.Floor(float64(cardID) / float64(rank)))
	var mask string
	switch rank {
	case model.Rank10000:
		mask = fmt.Sprintf("%dXXXX", normalized)
	case model.Rank1000:
		mask = fmt.Sprintf("%dXXX", normalized)
	case model.Rank100:
		mask = fmt.Sprintf("%dXX", normalized)
	default:
		log.Fatalf("[unity_repo] unknown rank %d", rank)
	}
	return mask
}
