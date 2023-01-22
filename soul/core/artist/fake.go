package artist

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"time"
)

type FakeArtist struct {
	cardRepository cardRepository
}

func NewFakeArtist(cardRepository cardRepository) *FakeArtist {
	return &FakeArtist{cardRepository: cardRepository}
}

func (fa *FakeArtist) GetCard(ctx context.Context, spell model.Spell) (model.Card, error) {
	time.Sleep(time.Second * 10)
	fakeNumber := rand.Intn(5)
	if b, err := os.ReadFile(fmt.Sprintf("files/fakes/%d.jpeg", fakeNumber)); err != nil {
		return model.Card{}, errors.Wrap(err, "[fake artist] failed to get file")
	} else {
		card := model.Card{
			Spell:   spell,
			Image:   b,
			Version: "v0.0",
		}
		return fa.cardRepository.SavePainting(ctx, card)
	}
}
