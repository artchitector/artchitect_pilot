package model

import "math"

func GetCardThousand(cardID uint) uint {
	return uint(math.Ceil(float64(cardID) / 10000))
}
