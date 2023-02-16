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
	SaveHundredImage(rank uint, hundred uint, data []byte) error
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

func (h *UploadHandler) HandleHundred(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msgf("[upload:hundred] failed to get file")
		c.String(http.StatusBadRequest, errors.Wrap(err, "[saver_upload] failed to get file from form").Error())
		return
	}
	rankStr := c.PostForm("rank")
	if rankStr == "" {
		log.Error().Msgf("[upload:hundred] rank must be integer")
		c.String(http.StatusBadRequest, "rank must be integer")
		return
	}
	rank, err := strconv.Atoi(rankStr)
	if err != nil {
		log.Error().Msgf("[upload:hundred] rank must be integer")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	hundredStr := c.PostForm("hundred")
	if hundredStr == "" {
		log.Error().Msgf("[upload:hundred] hundred must be integer")
		c.String(http.StatusBadRequest, "hundred must be integer")
		return
	}
	hundred, err := strconv.Atoi(hundredStr)
	if err != nil {
		log.Error().Msgf("[upload:hundred] hundred must be integer")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Info().Msgf("[saver_upload] incoming transmission for r:%d h:%d", rank, hundred)

	f, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msgf("[upload:hundred] failed to open file")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Error().Err(err).Msgf("[upload:hundred] failed to io.ReadAll")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.saver.SaveHundredImage(uint(rank), uint(hundred), data); err != nil {
		log.Error().Err(err).Msgf("[upload:hundred] failed to SaveHundredImage r:%d h:%d", rank, hundred)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Info().Msgf("[upload:HandleHundred] finished save r:%d h:%d", rank, hundred)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
