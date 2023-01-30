package engine

import (
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/pkg/errors"
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

	bts, err := io.ReadAll(response.Body)
	return bts, errors.Wrap(err, "[artist] failed to read response body")
}
