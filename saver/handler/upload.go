package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

type UploadRequest struct {
	Password string ``
}

type saver interface {
	SaveImage(cardID uint, data []byte) error
	SaveUnityImage(filename string, data []byte) error
}

type UploadHandler struct {
	saver saver
}

func NewUploadHandler(saver saver) *UploadHandler {
	return &UploadHandler{saver}
}

func (h *UploadHandler) Handle(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msgf("[upload:card] failed to get file")
		c.String(http.StatusBadRequest, errors.Wrap(err, "[saver_upload] failed to get file from form").Error())
		return
	}
	cardIdStr := c.PostForm("card_id")
	if cardIdStr == "" {
		log.Error().Msgf("[upload:card] card_id must be integer")
		c.String(http.StatusBadRequest, "card_id must be integer")
		return
	}
	cardID, err := strconv.Atoi(cardIdStr)
	if err != nil {
		log.Error().Msgf("[upload:card] card_id must be integer")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Info().Msgf("[saver_upload] incoming transmission for card id=%d", cardID)

	f, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msgf("[upload:card] failed to to open file")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Error().Err(err).Msgf("[upload:card] failed to io readAll")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.saver.SaveImage(uint(cardID), data); err != nil {
		log.Error().Err(err).Msgf("[upload:card] failed SaveImage card_id=%d", cardID)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

func (h *UploadHandler) HandleUnity(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msgf("[upload:unity] failed to get file")
		c.String(http.StatusBadRequest, errors.Wrap(err, "[saver_upload] failed to get file from form").Error())
		return
	}
	filename := c.PostForm("filename")
	if filename == "" {
		log.Error().Msgf("[upload:unity] filename must be")
		c.String(http.StatusBadRequest, "filename must be")
		return
	}

	log.Info().Msgf("[saver_upload] incoming unity transmission %s", filename)

	f, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msgf("[upload:unity] failed to open file")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Error().Err(err).Msgf("[upload:unity] failed to io.ReadAll")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.saver.SaveUnityImage(filename, data); err != nil {
		log.Error().Err(err).Msgf("[upload:hundred] failed to SaveHundredImage %s", filename)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Info().Msgf("[upload:HandleHundred] finished save %s", filename)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
