package handler

import (
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/artchitector/artchitect/resizer"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
)

var sizes = []string{model.SizeF, model.SizeM, model.SizeS, model.SizeXS}

type UploadRequest struct {
	Password string ``
}

type UploadHandler struct {
	cardsPath string
}

func NewUploadHandler(cardsPath string) *UploadHandler {
	return &UploadHandler{cardsPath}
}

func (h *UploadHandler) Handle(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, errors.Wrap(err, "[saver_upload] failed to get file from form").Error())
		return
	}
	cardIdStr := c.PostForm("card_id")
	if cardIdStr == "" {
		c.String(http.StatusBadRequest, "card_id must be integer")
		return
	}
	cardID, err := strconv.Atoi(cardIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Info().Msgf("[saver_upload] incoming transmission for card id=%d", cardID)

	f, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	if err := h.saveUploadedFile(data, uint(cardID)); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
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
func (h *UploadHandler) saveUploadedFile(data []byte, cardID uint) error {
	for _, size := range sizes {
		resized, err := resizer.ResizeBytes(data, size)
		if err != nil {
			return errors.Wrapf(err, "[saver_upload] failed to resize card %d, %s", cardID, size)
		}

		idFolder := fmt.Sprintf("%d", int(math.Ceil(float64(cardID)/10000)))
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
