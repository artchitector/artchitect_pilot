package resizer

import (
	"bytes"
	"github.com/artchitector/artchitect/model"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"image"
	"image/jpeg"
	"math"
	"time"
)

func ResizeImage(img image.Image, size string) (image.Image, error) {
	dimension := float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
	var height, width uint

	switch size {
	case model.SizeXF:
		// image already full size, but downgrade quality
		// it's 2560x3840 (it's 5Mb size with 100% quality, down to 90%)
		width = uint(img.Bounds().Dx())
		height = uint(math.Round(float64(width) * dimension))
	case model.SizeF:
		width = uint(1024)
		height = uint(math.Round(float64(width) * dimension))
	case model.SizeM:
		width = uint(512)
		height = uint(math.Round(float64(width) * dimension))
	case model.SizeS:
		width = uint(256)
		height = uint(math.Round(float64(width) * dimension))
	case model.SizeXS:
		width = uint(128)
		height = uint(math.Round(float64(width) * dimension))
	default:
		// TODO сделать из этого ответ bad-requst, если такое пришло
		return nil, errors.Errorf("[resizer] wrong size %s", size)
	}

	img = resize.Resize(width, height, img, resize.Lanczos3)
	return img, nil
}

func ResizeBytes(rawImg []byte, size string) ([]byte, error) {
	var quality int
	switch size {
	case model.SizeXF:
		quality = model.QualityXF
	case model.SizeF:
		quality = model.QualityF
	case model.SizeM:
		quality = model.QualityM
	case model.SizeS:
		quality = model.QualityS
	case model.SizeXS:
		quality = model.QualityXS
	}
	return ResizeBytesWithQuality(rawImg, size, quality)
}

func ResizeBytesWithQuality(rawImg []byte, size string, quality int) ([]byte, error) {
	start := time.Now()
	r := bytes.NewReader(rawImg)
	img, err := jpeg.Decode(r)
	if err != nil {
		return []byte{}, errors.Wrap(err, "[card_handler] failed to decode jpeg")
	}

	img, err = ResizeImage(img, size)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "[resizer] failed to resize image %s/%d", size, quality)
	}

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
