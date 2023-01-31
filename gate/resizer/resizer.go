package resizer

import (
	"bytes"
	"github.com/artchitector/artchitect/model"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image/jpeg"
	"math"
	"time"
)

func Resize(rawImg []byte, size string) ([]byte, error) {
	start := time.Now()
	r := bytes.NewReader(rawImg)
	img, err := jpeg.Decode(r)
	if err != nil {
		return []byte{}, errors.Wrap(err, "[card_handler] failed to decode jpeg")
	}

	dimension := float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
	var height, width uint
	var quality int
	switch size {
	case model.SizeXF:
		// nothing to do, image already full
		// it's 2560x3840 (it's 5Mb size, and it's not cached in redis)
		return rawImg, nil
	case model.SizeF:
		width = uint(1024)
		height = uint(math.Round(float64(width) * dimension))
		quality = 90
	case model.SizeM:
		width = uint(512)
		height = uint(math.Round(float64(width) * dimension))
		quality = 80
	case model.SizeS:
		width = uint(256)
		height = uint(math.Round(float64(width) * dimension))
		quality = 75
	case model.SizeXS:
		width = uint(128)
		height = uint(math.Round(float64(width) * dimension))
		quality = 75
	default:
		// TODO сделать из этого ответ bad-requst, если такое пришло
		return []byte{}, errors.Errorf("[resizer] wrong size %s", size)
	}

	img = resize.Resize(width, height, img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return []byte{}, errors.Wrapf(err, "[resizer] failed to encode jpeg")
	}
	log.Info().Msgf(
		"[resizer] resized, size=%s, time: %s",
		size,
		time.Now().Sub(start),
	)

	return buf.Bytes(), nil
}
