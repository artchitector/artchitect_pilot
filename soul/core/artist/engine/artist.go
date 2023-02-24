package engine

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
	"image"
	"image/png"
	"io"
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

func (e *ArtistEngine) GetImage(ctx context.Context, spell model.Spell) (image.Image, error) {
	client := http.Client{
		Timeout: time.Second * 90,
	}
	response, err := client.PostForm(e.artistURL+"/painting", url.Values{
		"tags":    {spell.Tags},
		"seed":    {fmt.Sprintf("%d", spell.Seed)},
		"width":   {"640"},
		"height":  {"960"},
		"steps":   {"50"},
		"upscale": {"4"},
		"version": {spell.Version},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to artist")
	}
	defer response.Body.Close()

	bts, err := io.ReadAll(response.Body)
	r := bytes.NewReader(bts)
	img, err := png.Decode(r)
	if err != nil {
		return nil, errors.Wrap(err, "[artist] failed to get valid jpeg")
	}

	return img, errors.Wrap(err, "[artist] failed to read response body")
}
