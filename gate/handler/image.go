package handler

import (
	"fmt"
	mmrPkg "github.com/artchitector/artchitect/memory"
	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"net/http"
	"os"
)

type ImageRequest struct {
	ID   uint   `uri:"id" binding:"numeric"`
	Size string `uri:"size" binding:"required"` // size f - full, size m - 2-times smaller dimensions, size s - 4-times smaller dimensions
}

type ImageUnityRequest struct {
	Mask    string `uri:"mask" binding:"required"`
	Version string `uri:"version" binding:"required"`
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

	if request.ID == 0 {
		if dt, err := os.ReadFile(fmt.Sprintf("./files/black-%s.jpg", request.Size)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.Data(http.StatusOK, "image/jpeg", dt)
			return
		}
	}

	imageBytes, err := ih.memory.GetCardImage(c, request.ID, request.Size)
	if err != nil {
		log.Error().Err(err).Msgf("[image_handler] failed to GetCardImage")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", imageBytes)
}

func (ih *ImageHandler) HandleUnity(c *gin.Context) {
	var request ImageUnityRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Size == model.SizeXF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "XF size is not supported"})
		return
	}

	imgBytes, err := ih.memory.GetUnityImage(c, request.Mask, request.Size, request.Version)
	if err != nil {
		if errors.Is(err, mmrPkg.ErrNotFound) {
			if dt, err := os.ReadFile(fmt.Sprintf("./files/black-%s.jpg", request.Size)); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				c.Data(http.StatusOK, "image/jpeg", dt)
				return
			}
		}
		log.Error().Err(err).Msgf("[image_handler] failed to GetHundredImage %s %s", request.Mask, request.Size)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", imgBytes)
}
