package unity_worker

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"math"
)

type cardsRepo interface {
	GetMaxCardID(ctx context.Context) (uint, error)
}

type unityRepo interface {
	GetUnity(mask string) (model.Unity, error)
	CreateUnity(mask string) (model.Unity, error)
}

type UnityWorker struct {
	cardsRepo cardsRepo
	unityRepo unityRepo
}

func NewUnityWorker(cardsRepo cardsRepo, unityRepo unityRepo) *UnityWorker {
	return &UnityWorker{cardsRepo: cardsRepo, unityRepo: unityRepo}
}

func (u *UnityWorker) Work(ctx context.Context) {
	max, err := u.cardsRepo.GetMaxCardID(ctx)
	if err != nil {
		log.Error().Err(err).Msgf("[unity_worker] failed to get max card id")
		return
	}
	log.Info().Msgf("[unity_worker] got max %d", max)
	n := math.Ceil(float64(max) / model.Rank10000)
	for i := 0; i < int(n); i++ {
		mask := fmt.Sprintf("%dXXXX", i)
		_, err := u.unityRepo.GetUnity(mask)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Err(err).Msgf("[unity_worker] failed to get unity %s", mask)
			return
		} else if err == nil {
			log.Info().Msgf("[unity_worker] unity %s already exists", mask)
			continue // already exists
		}
		if _, err := u.unityRepo.CreateUnity(mask); err != nil {
			log.Error().Err(err).Msgf("[unity_worker] failed create unity %s", mask)
			return
		}
		log.Info().Msgf("[unity_worker] created unity %s", mask)
	}
}
