package artist

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"time"
)

type FakeArtist struct {
	cardRepository cardRepository
	notifier       notifier
}

func NewFakeArtist(cardRepository cardRepository, notifier notifier) *FakeArtist {
	return &FakeArtist{cardRepository, notifier}
}

func (fa *FakeArtist) GetCard(ctx context.Context, spell model.Spell, artistState *model.ArtistState) (model.Card, error) {
	lastPaintingTime, err := fa.cardRepository.GetLastCardPaintTime(ctx)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[artist] failed to get LastPaintingTime")
	}

	paintStart := time.Now()
	updaterCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case <-updaterCtx.Done():
				return
			case <-time.NewTicker(time.Millisecond * 1000).C:
				artistState.LastCardPaintTime = lastPaintingTime
				artistState.CurrentCardPaintTime = uint64(time.Now().Sub(paintStart).Seconds())
				if err := fa.notifier.NotifyArtistState(ctx, *artistState); err != nil {
					log.Error().Err(err).Msg("[artist] failed to notift artist state")
				}
			}
		}
	}()

	time.Sleep(time.Second * 10)
	fakeNumber := rand.Intn(5)

	if b, err := os.ReadFile(fmt.Sprintf("files/fakes/%d.jpeg", fakeNumber)); err != nil {
		return model.Card{}, errors.Wrap(err, "[fake artist] failed to get file")
	} else {
		cancel()
		card := model.Card{
			Spell:     spell,
			Image:     b,
			Version:   "v0.0",
			PaintTime: uint64(time.Now().Sub(paintStart).Seconds()),
		}
		return fa.cardRepository.SavePainting(ctx, card)
	}
}
