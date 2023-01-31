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
	var quality int
	switch size {
	case model.SizeXF:
		// nothing to do, image already full
		return rawImg, nil
	case model.SizeF:
		width = uint(1024)
		height = uint(1536)
		quality = 90
	case model.SizeM:
		width = uint(512)
		height = uint(768)
		quality = 80
	case model.SizeS:
		width = uint(256)
		height = uint(384)
		quality = 80
	case model.SizeXS:
		width = uint(128)
		height = uint(192)
		quality = 80
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
