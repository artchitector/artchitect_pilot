package fake

import (
	"context"
	"github.com/pkg/errors"
)

type FakeOrigin struct {
}

func (o *FakeOrigin) Select(ctx context.Context, totalVariants uint) (uint, error) {
	return 0, errors.Errorf("fake origin disabled")
}
