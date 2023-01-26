package engine

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"image/jpeg"
	"net/http"
	"net/url"
	"time"
)

type ArtistEngine struct {
	artistURL string
}

func NewArtistEngine(artistURL string) *ArtistEngine {
	return &ArtistEngine{artistURL}
}

func (e *ArtistEngine) GetImage(ctx context.Context, spell model.Spell) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 90,
	}
	response, err := client.PostForm(e.artistURL+"/painting", url.Values{
		"tags": {spell.Tags},
		"seed": {fmt.Sprintf("%d", spell.Seed)},
	})
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to make request to artist")
	}
	defer response.Body.Close()

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to decode jpeg from response")
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return []byte{}, errors.Wrap(err, "failed to encode image into jpeg data")
	}

	return buf.Bytes(), nil
}
