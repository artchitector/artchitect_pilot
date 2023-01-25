package artist

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image/jpeg"
	"net/http"
	"net/url"
	"time"
)

type notifier interface {
	NotifyArtistState(ctx context.Context, state model.ArtistState) error
}

type cardRepository interface {
	SavePainting(ctx context.Context, painting model.Card) (model.Card, error)
	GetLastCardPaintTime(ctx context.Context) (uint64, error)
}

type Artist struct {
	artistURL string
	cardRepo  cardRepository
	notifier  notifier
}

func NewArtist(artistURL string, paintingRepository cardRepository, notifier notifier) *Artist {
	return &Artist{artistURL, paintingRepository, notifier}
}

func (a *Artist) GetCard(ctx context.Context, spell model.Spell, artistState *model.ArtistState) (model.Card, error) {
	lastPaintingTime, err := a.cardRepo.GetLastCardPaintTime(ctx)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[artist] failed to get LastPaintingTime")
	}
	client := http.Client{
		Timeout: time.Second * 90,
	}
	log.Info().Msgf("Start get painting process from artist. tags: %s, seed: %d", spell.Tags, spell.Seed)
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
				if err := a.notifier.NotifyArtistState(ctx, *artistState); err != nil {
					log.Error().Err(err).Msg("[artist] failed to notift artist state")
				}
			}
		}
	}()

	response, err := client.PostForm(a.artistURL+"/painting", url.Values{
		"tags": {spell.Tags},
		"seed": {fmt.Sprintf("%d", spell.Seed)},
	})
	if err != nil {
		return model.Card{}, errors.Wrap(err, "failed to make request to artist")
	}
	defer response.Body.Close()
	cancel()

	paintTime := time.Now().Sub(paintStart)

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "failed to decode jpeg from response")
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return model.Card{}, errors.Wrap(err, "failed to encode image into jpeg data")
	}
	painting := model.Card{
		Spell:     spell,
		Image:     buf.Bytes(),
		Version:   spell.Version,
		PaintTime: uint64(paintTime.Seconds()),
	}
	painting, err = a.cardRepo.SavePainting(ctx, painting)
	log.Info().Msgf("Received and saved painting from artist: id=%d", painting.ID)
	return painting, err
}
