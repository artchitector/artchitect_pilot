package engine

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"time"
)

type FakeEngine struct {
	fakeGenerationTime uint
}

func NewFakeEngine(fakeGenerationTime uint) *FakeEngine {
	return &FakeEngine{
		fakeGenerationTime,
	}
}

func (e *FakeEngine) GetImage(ctx context.Context, spell model.Spell) (image.Image, error) {
	fakeNumber := rand.Intn(20) + 1
	if b, err := os.ReadFile(fmt.Sprintf("files/fakes/%d.jpeg", fakeNumber)); err != nil {
		return nil, errors.Wrap(err, "[fake artist] failed to get file")
	} else {
		time.Sleep(time.Second * time.Duration(e.fakeGenerationTime)) // imitation of long-running process
		buf := bytes.NewBuffer(b)
		img, err := jpeg.Decode(buf)
		return img, errors.Wrap(err, "[fake_artist] failed to decode jpeg")
	}
}
