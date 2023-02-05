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

	if err := h.saver.SaveImage(uint(cardID), data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
