package driver

import (
	"bytes"
	"context"
	"crypto/md5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
)

/*
WebcamDriver gets image from webcam (from http-url, where you can get jpg-image from webcam).
Then image normalized into float64 value from 0.0 to 1.0. And this is the answer for all questions.
*/
type WebcamDriver struct {
	originUrl string
}

func NewWebcamDriver(originUrl string) *WebcamDriver {
	return &WebcamDriver{originUrl}
}

func (w *WebcamDriver) GetValue(ctx context.Context) (float64, error) {
	response, err := http.Get(w.originUrl)
	if err != nil {
		return 0.0, errors.Wrapf(err, "failed to get %s", w.originUrl)
	}
	defer response.Body.Close()

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		return 0.0, errors.Wrap(err, "failed to decode image from response.Body")
	}

	result, err := w.imageToNumber(ctx, img)
	log.Info().Msgf("[webcam] got number %f", result)
	return result, err
}

func (w *WebcamDriver) imageToNumber(ctx context.Context, originalImg image.Image) (float64, error) {
	var result float64
	var err error
	result, err = w.imageToNumberHash(ctx, originalImg)
	if err != nil {
		return 0.0, errors.Wrapf(err, "[webcam] failed to imageToNumberHash")
	}

	return result, nil
}

func (w *WebcamDriver) imageToNumberHash(ctx context.Context, originalImg image.Image) (float64, error) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, originalImg); err != nil {
		return 0.0, errors.Wrap(err, "failed to encode png to bytes")
	}

	hash := md5.Sum(buf.Bytes())

	var result uint
	for _, b := range hash[:8] {
		result = (result << 8) | uint(b)
	}

	flResult := float64(result) / float64(math.MaxUint)
	log.Debug().Msgf("[webcam][imageToNumberHash] generated number: %d. Meaning: %.12f", result, flResult)
	return flResult, nil
}
