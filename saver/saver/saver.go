package saver

import (
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/artchitector/artchitect/resizer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"path"
)

var sizes = []string{model.SizeF, model.SizeM, model.SizeS, model.SizeXS}

type Saver struct {
	cardsPath    string
	unityPath    string
	fullsizePath string
}

func NewSaver(cardsPath string, unityPath string, fullsizePath string) *Saver {
	return &Saver{cardsPath, unityPath, fullsizePath}
}

/*
file structure:
  - all images are in /var/artchitect/arts folder (set in env)
  - 10k-arts is in separate folder folder=(id % 10000)
  - arts names in these folders:
    art-56910-f.jpg
    art-56910-m.jpg
    art-56910-s.jpg
    art-56910-xs.jpg
*/
func (h *Saver) SaveImage(cardID uint, data []byte) error {
	for _, size := range sizes {
		resized, err := resizer.ResizeBytes(data, size)
		if err != nil {
			return errors.Wrapf(err, "[saver_upload] failed to resize card %d, %s", cardID, size)
		}

		idFolder := fmt.Sprintf("%d", model.GetCardThousand(cardID))
		folderPath := path.Join(h.cardsPath, idFolder)
		filename := fmt.Sprintf("art-%d-%s.jpg", cardID, size)

		err = h.saveFile(
			folderPath,
			filename,
			resized,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Saver) SaveUnityImage(unityName string, data []byte) error {
	for _, size := range sizes {
		resized, err := resizer.ResizeBytes(data, size)
		if err != nil {
			return errors.Wrapf(err, "[saver_upload] failed to resize %s s:%s", unityName, size)
		}

		err = h.saveFile(
			h.unityPath,
			fmt.Sprintf("unity-%s-%s.jpg", unityName, size),
			resized,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Saver) SaveFullsizeArt(cardID uint, data []byte) error {
	idFolder := fmt.Sprintf("%d", model.GetCardThousand(cardID))
	folderPath := path.Join(h.fullsizePath, idFolder)
	filename := fmt.Sprintf("art-%d.jpg", cardID)

	err := h.saveFile(
		folderPath,
		filename,
		data,
	)
	return err
}

func (h *Saver) saveFile(folder string, filename string, data []byte) error {
	folderPath := path.Join(folder)
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "[saver_upload] failed to create folder %s", folderPath)
	}

	p := path.Join(folderPath, filename)
	err := os.WriteFile(p, data, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "[saver] failed to save file %s", p)
	} else {
		log.Info().Msgf("[saver] saved file %s. size=%d", p, len(data))
	}
	return nil
}
