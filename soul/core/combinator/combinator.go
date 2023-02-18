package combinator

import (
	"bytes"
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/artchitector/artchitect/resizer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"math"
)

type cardRepository interface {
	GetAnyCardIDFromHundred(ctx context.Context, rank uint, hundred uint) (uint, error)
}

type memory interface {
	DownloadImage(ctx context.Context, cardID uint, size string) ([]byte, error)
}

type saver interface {
	SaveUnity(filename string, imgFile []byte) error
}

type hundredRepository interface {
	SaveHundred(rank uint, hundred uint) (model.Hundred, error)
}

type watermark interface {
	AddUnityWatermark(originalImage image.Image, mask string) (image.Image, error)
}

type Combinator struct {
	cardRepository    cardRepository
	memory            memory
	saver             saver
	hundredRepository hundredRepository
	watermark         watermark
}

func NewCombinator(cardRepository cardRepository, memory memory, saver saver, hundredRepository hundredRepository, watermark watermark) *Combinator {
	return &Combinator{cardRepository, memory, saver, hundredRepository, watermark}
}

func (c *Combinator) CombineThumb(ctx context.Context, cardIDs []uint, mask string) error {
	var imgs []image.Image
	size := model.SizeS
	if len(cardIDs) > 16 {
		size = model.SizeXS
	}
	for _, cardID := range cardIDs {
		imageFile, err := c.memory.DownloadImage(ctx, cardID, size)
		if err != nil {
			return errors.Wrapf(err, "[combinator] failed to get image %d %s", cardID, size)
		} else {
			r := bytes.NewReader(imageFile)
			img, err := jpeg.Decode(r)
			if err != nil {
				return errors.Wrapf(err, "[combinator] failed to decode jpeg %d %s", cardID, model.SizeM)
			}
			imgs = append(imgs, img)
		}
	}

	thumb := c.combineTotal(imgs)
	resizedThumb, err := resizer.ResizeImage(thumb, model.SizeF)
	if err != nil {
		return errors.Wrapf(err, "failed to resize total image")
	}

	resizedThumb, err = c.watermark.AddUnityWatermark(resizedThumb, mask)
	if err != nil {
		return errors.Wrapf(err, "failed to resize add watermark")
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, resizedThumb, &jpeg.Options{Quality: model.QualityF}); err != nil {
		return errors.Wrapf(err, "failed to encode total jpeg image")
	}
	err = c.saver.SaveUnity(mask, buf.Bytes())
	if err != nil {
		return errors.Wrapf(err, "[combinator] failed to combine total card")
	}
	log.Info().Msgf("[combinator] combined thumb from %+v", cardIDs)
	return nil
}

func (c *Combinator) combineTotal(imgs []image.Image) image.Image {
	size := int(math.Sqrt(float64(len(imgs))))
	first := imgs[0]
	total := image.NewRGBA(image.Rect(0, 0, first.Bounds().Dx()*size, first.Bounds().Dy()*size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			idx := y*size + x
			pnt := image.Point{-x * imgs[idx].Bounds().Dx(), -y * imgs[idx].Bounds().Dy()}
			draw.Draw(total, total.Bounds(), imgs[idx], pnt, draw.Over)
			log.Info().Msgf("[combinator] draw over x:%d y:%d idx:%d point:%+v", x, y, idx, pnt)
		}
	}
	return total
}
