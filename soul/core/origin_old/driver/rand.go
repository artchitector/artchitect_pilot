package driver

import (
	"context"
	"github.com/rs/zerolog/log"
	"math/rand"
)

type RandDriver struct {
}

func NewRandDriver() *RandDriver {
	return &RandDriver{}
}

func (r RandDriver) GetValue(ctx context.Context) (float64, error) {
	v := float64(rand.Int63n(1000000)) / 1000000.0
	log.Info().Msgf("[rand_driver] val=%f", v)
	return v, nil
}
