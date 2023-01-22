package resizer

import (
	"bytes"
	"github.com/artchitector/artchitect/model"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image/jpeg"
	"time"
)

func Resize(rawImg []byte, size string) ([]byte, error) {
	start := time.Now()
	r := bytes.NewReader(rawImg)
	img, err := jpeg.Decode(r)
	if err != nil {
		return []byte{}, errors.Wrap(err, "[card_handler] failed to decode jpeg")
	}

	var height, width uint
	switch size {
	case model.SizeF:
		// nothing to do, image already full
		return rawImg, nil
	case model.SizeM:
		width = uint(img.Bounds().Size().X / 2)
		height = uint(img.Bounds().Size().Y / 2)
	case model.SizeS:
		width = uint(img.Bounds().Size().X / 4)
		height = uint(img.Bounds().Size().Y / 4)
	case model.SizeXS:
		width = uint(img.Bounds().Size().X / 8)
		height = uint(img.Bounds().Size().Y / 8)
	default:
		// TODO сделать из этого ответ bad-requst, если такое пришло
		return []byte{}, errors.Errorf("[resizer] wrong size %s", size)
	}

	img = resize.Resize(width, height, img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		return []byte{}, errors.Wrapf(err, "[resizer] failed to encode jpeg")
	}
	log.Info().Msgf(
		"[resizer] resized, size=%s, time: %s",
		size,
		time.Now().Sub(start),
	)

	return buf.Bytes(), nil
}
