package saver

import (
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/artchitector/artchitect/resizer"
	"github.com/pkg/errors"
	"os"
	"path"
)

var sizes = []string{model.SizeF, model.SizeM, model.SizeS, model.SizeXS}

type Saver struct {
	cardsPath string
}

func NewSaver(cardsPath string) *Saver {
	return &Saver{cardsPath}
}

/*
file structure:
  - all images are in /root/cards folder (specified from .env variable)
  - every 10k cards is in separate folder: folder=(id % 10000)
  - card names in these folders:
    card-56910-f.jpg
    card-56910-m.jpg
    card-56910-s.jpg
    card-56910-xs.jpg
    these files statically served by nginx, and gate services can take img and proxy it
*/
func (h *Saver) SaveImage(cardID uint, data []byte) error {
	for _, size := range sizes {
		resized, err := resizer.ResizeBytes(data, size)
		if err != nil {
			return errors.Wrapf(err, "[saver_upload] failed to resize card %d, %s", cardID, size)
		}

		idFolder := fmt.Sprintf("%d", model.GetCardThousand(cardID))
		folderPath := path.Join(h.cardsPath, idFolder)
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return errors.Wrapf(err, "[saver_upload] failed to create folder")
		}

		p := path.Join(folderPath, fmt.Sprintf("card-%d-%s.jpg", cardID, size))
		err = os.WriteFile(p, resized, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "[saver_upload] failed to save file %s", p)
		}
	}
	return nil
}
