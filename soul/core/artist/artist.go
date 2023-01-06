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

type paintingRepository interface {
	SavePainting(ctx context.Context, painting model.Painting) (model.Painting, error)
}

type Artist struct {
	artistURL          string
	paintingRepository paintingRepository
}

func NewArtist(artistURL string, paintingRepository paintingRepository) *Artist {
	return &Artist{artistURL, paintingRepository}
}

func (a *Artist) GetPainting(ctx context.Context, spell model.Spell) (model.Painting, error) {
	client := http.Client{
		Timeout: time.Second * 90,
	}
	log.Info().Msgf("Start get painting process from artist. tags: %s, seed: %d", spell.Tags, spell.Seed)
	response, err := client.PostForm(a.artistURL+"/painting", url.Values{
		"tags": {spell.Tags},
		"seed": {fmt.Sprintf("%d", spell.Seed)},
	})
	if err != nil {
		return model.Painting{}, errors.Wrap(err, "failed to make request to artist")
	}
	defer response.Body.Close()

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		return model.Painting{}, errors.Wrap(err, "failed to decode jpeg from response")
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return model.Painting{}, errors.Wrap(err, "failed to encode image into jpeg data")
	}
	painting := model.Painting{
		Spell: spell,
		Image: buf.Bytes(),
	}
	painting, err = a.paintingRepository.SavePainting(ctx, painting)
	log.Info().Msgf("Received and saved painting from artist: id=%d", painting.ID)
	return painting, err
}
