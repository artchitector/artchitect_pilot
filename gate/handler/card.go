package handler

import (
	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
)

type CardRequest struct {
	ID uint `uri:"id" binding:"required,numeric"`
}

type ImageRequest struct {
	ID   uint   `uri:"id" binding:"required,numeric"`
	Size string `uri:"size" binding:"required"` // size f - full, size m - 2-times smaller dimensions, size s - 4-times smaller dimensions
}

type CardHandler struct {
	cardsRepository cardsRepository
	cache           cache
	memory          memory
}

func NewCardHandler(cardsRepository cardsRepository, cache cache, memory memory) *CardHandler {
	return &CardHandler{cardsRepository, cache, memory}
}

func (lh *CardHandler) Handle(c *gin.Context) {
	var request CardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card, err := lh.cache.GetCard(c, uint(request.ID))
	if err != nil {
		log.Error().Err(err).Msgf("[card_handler:Handle] failed to get card(id=%d) from cache", request.ID)
	} else {
		c.JSON(http.StatusOK, card)
		return
	}

	card, err = lh.cardsRepository.GetCard(c, request.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}

func (ch *CardHandler) HandleImage(c *gin.Context) {
	var request ImageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Size == model.SizeXF {
		c.JSON(http.StatusBadRequest, gin.H{"error": "XF size is not supported in gate"})
		return
	}

	imageBytes, err := ch.memory.GetImage(c, request.ID, request.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", imageBytes)
}
