package handler

import (
	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type ImageRequest struct {
	ID   uint   `uri:"id" binding:"required,numeric"`
	Size string `uri:"size" binding:"required"` // size f - full, size m - 2-times smaller dimensions, size s - 4-times smaller dimensions
}

type ImageHundredRequest struct {
	Rank    uint   `uri:"rank" binding:"required,numeric"`
	Hundred uint   `uri:"hundred" binding:"numeric"`
	Size    string `uri:"size" binding:"required"`
}

type ImageHandler struct {
	memory memory
}

func NewImageHandler(memory memory) *ImageHandler {
	return &ImageHandler{memory}
}

func (ih *ImageHandler) HandleImage(c *gin.Context) {
	var request ImageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Size == model.SizeXF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "XF size is not supported"})
		return
	}

	imageBytes, err := ih.memory.GetCardImage(c, request.ID, request.Size)
	if err != nil {
		log.Error().Err(err).Msgf("[image_handler] failed to GetCardImage")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", imageBytes)
}

func (ih *ImageHandler) HandleHundred(c *gin.Context) {
	var request ImageHundredRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Size == model.SizeXF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "XF size is not supported"})
		return
	}

	imgBytes, err := ih.memory.GetHundredImage(c, request.Rank, request.Hundred, request.Size)
	if err != nil {
		log.Error().Err(err).Msgf("[image_handler] failed to GetHundredImage r:%d h:%d", request.Rank, request.Hundred)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", imgBytes)
}
