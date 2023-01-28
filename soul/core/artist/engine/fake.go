package engine

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"time"
)

type FakeEngine struct {
}

func NewFakeEngine() *FakeEngine {
	return &FakeEngine{}
}

func (e *FakeEngine) GetImage(ctx context.Context, spell model.Spell) ([]byte, error) {
	fakeNumber := rand.Intn(20) + 1
	if b, err := os.ReadFile(fmt.Sprintf("files/fakes/%d.jpeg", fakeNumber)); err != nil {
		return []byte{}, errors.Wrap(err, "[fake artist] failed to get file")
	} else {
		time.Sleep(time.Second * 10) // imitation of long-running process
		return b, nil
	}
}
