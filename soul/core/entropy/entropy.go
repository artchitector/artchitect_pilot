package entropy

import (
	"context"
	"math"
)

type Entropy struct {
	lightmaster *Lightmaster
}

func NewEntropy(lightmaster *Lightmaster) *Entropy {
	return &Entropy{lightmaster: lightmaster}
}

/*
	Artchitect asks "select one element from set, i have total 100 elements.
	Entropy replies: "take element 31" (calculated with the lightnoise-entropy)
*/

func (e *Entropy) Select(ctx context.Context, totalElements uint) (uint, error) {
	entropyF := e.lightmaster.GetChoice(ctx)
	targetIndex := uint(math.Floor(float64(totalElements) * entropyF))

	return targetIndex, nil
}
