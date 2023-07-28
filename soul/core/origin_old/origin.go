package origin_old

import (
	"context"
	"github.com/pkg/errors"
	"math"
)

// Origin - source of everything. Origin generates random seed based on God's will.
// Here I use web-camera, grab picture from webcam, then I translate this "quantum noise" into random number.

// driver is randomNumberGetter interface. We can take new random value with GetValue method.
type driver interface {
	// GetValue returns float64 from 0 to 1
	GetValue(ctx context.Context) (float64, error)
}

type Origin struct {
	provider driver
}

func NewOrigin(provider driver) *Origin {
	return &Origin{provider}
}

func (o *Origin) Select(ctx context.Context, totalVariants uint) (uint, error) {
	val, err := o.provider.GetValue(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "[origin] failed to getValue from provider")
	}
	return uint(math.Round(float64(totalVariants-1) * val)), nil
}
