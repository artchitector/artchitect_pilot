package resources

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"net/http"
)

type Webcam struct {
	originUrl string
}

func (w *Webcam) GetStream(ctx context.Context) chan image.Image {
	ch := make(chan image.Image)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("[webcam] stop reading stream")
			default:
				if img, err := w.getFrame(ctx); err != nil {
					log.Error().Err(err).Msgf("[webcam] failed getFrame")
				} else {
					log.Info().Msgf("[webcam] got image")
					ch <- img
				}
			}
		}
	}()
	return ch
}

func (w *Webcam) getFrame(ctx context.Context) (image.Image, error) {
	response, err := http.Get(w.originUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s", w.originUrl)
	}
	defer response.Body.Close()

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image from response.Body")
	}
	img = w.yCrCb2RGBA(img)
	return img, nil
}

func (w *Webcam) yCrCb2RGBA(img image.Image) image.Image {
	b := img.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)
	return m
}
