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
)

const (
	Width  = 4
	Height = 4
)

type cardRepository interface {
	GetAnyCardIDFromHundred(ctx context.Context, rank uint, hundred uint) (uint, error)
}

type memory interface {
	DownloadImage(ctx context.Context, cardID uint, size string) ([]byte, error)
}

type saver interface {
	SaveHundred(rank uint, hundred uint, imgFile []byte) error
}

type hundredRepository interface {
	SaveHundred(rank uint, hundred uint) (model.Hundred, error)
}

type Combinator struct {
	cardRepository    cardRepository
	memory            memory
	saver             saver
	hundredRepository hundredRepository
}

func NewCombinator(cardRepository cardRepository, memory memory, saver saver, hundredRepository hundredRepository) *Combinator {
	return &Combinator{cardRepository, memory, saver, hundredRepository}
}

// CombineHundred - combines image matrix from all hundred (take 16 any images and make collage)
func (c *Combinator) CombineHundred(ctx context.Context, rank uint, hundred uint) error {
	var imgs []image.Image
	totalImages := Width * Height
	for i := 0; i < totalImages; i++ {
		cardID, err := c.cardRepository.GetAnyCardIDFromHundred(ctx, rank, hundred)
		if err != nil {
			return errors.Wrapf(err, "[combinator] failed to get card from r:%d h:%d", rank, hundred)
		}
		log.Info().Msgf("[combinator] selected card %d for r:%d h:%d", cardID, rank, hundred)
		imageFile, err := c.memory.DownloadImage(ctx, cardID, model.SizeS)
		if err != nil {
			return errors.Wrapf(err, "[combinator] failed to get image %d %s", cardID, model.SizeM)
		}
		r := bytes.NewReader(imageFile)
		img, err := jpeg.Decode(r)
		if err != nil {
			return errors.Wrapf(err, "[combinator] failed to decode jpeg %d %s", cardID, model.SizeM)
		}
		imgs = append(imgs, img)
	}

	totalImg := c.combineTotal(imgs)
	resizedTotal, err := resizer.ResizeImage(totalImg, model.SizeF)
	if err != nil {
		return errors.Wrapf(err, "failed to resize total image r:%d h:%d", rank, hundred)
	}
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, resizedTotal, &jpeg.Options{Quality: model.QualityF}); err != nil {
		return errors.Wrapf(err, "failed to encode total jpeg image r:%d h:%h")
	}
	err = c.saver.SaveHundred(rank, hundred, buf.Bytes())
	if err != nil {
		return errors.Wrapf(err, "[combinator] failed to combine total card")
	}
	h, err := c.hundredRepository.SaveHundred(rank, hundred)
	if err != nil {
		return errors.Wrapf(err, "[combinator] failed to save hundred r:%d h:%d", rank, hundred)
	}
	log.Info().Msgf("[combinator] Saved hundred r:%d h:%d id=%d", rank, hundred, h.ID)
	return nil
}

func (c *Combinator) combineTotal(imgs []image.Image) image.Image {
	first := imgs[0]
	total := image.NewRGBA(image.Rect(0, 0, first.Bounds().Dx()*Width, first.Bounds().Dy()*Height))
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			idx := y*Width + x
			pnt := image.Point{-x * imgs[idx].Bounds().Dx(), -y * imgs[idx].Bounds().Dy()}
			draw.Draw(total, total.Bounds(), imgs[idx], pnt, draw.Over)
			log.Info().Msgf("[combinator] draw over x:%d y:%d idx:%d point:%+v", x, y, idx, pnt)
		}
	}
	return total
}
