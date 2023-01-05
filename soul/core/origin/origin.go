package origin

import (
	"context"
	"github.com/pkg/errors"
	"math"
)

// Origin - source of everything. Origin generates random seed based on God's will.
// Here I use web-camera, grab picture from webcam, then I translate this "quantum noise" into random number.

type Origin struct {
	provider Driver
}

func NewOrigin(provider Driver) *Origin {
	return &Origin{provider}
}

// Driver is randomNumberGetter interface. We can take new random value with GetValue method. Min and Max made to understand whole scale.
type Driver interface {
	// GetValue returns float64 from 0 to 1
	GetValue(ctx context.Context) (float64, error)
}

func (o *Origin) YesNo(ctx context.Context) (bool, error) {
	val, err := o.provider.GetValue(ctx)
	if err != nil {
		return false, errors.Wrap(err, "[origin] failed to getValue from provider")
	}
	return val > 0.5, nil
}

func (o *Origin) Select(ctx context.Context, totalVariants uint64) (uint64, error) {
	val, err := o.provider.GetValue(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "[origin] failed to getValue from provider")
	}
	return uint64(math.Round(float64(totalVariants) * val)), nil
}
