package handler

import (
	"github.com/artchitector/artchitect/gate/resizer"
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
}

func NewCardHandler(cardsRepository cardsRepository, cache cache) *CardHandler {
	return &CardHandler{cardsRepository, cache}
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

	cached, err := ch.cache.GetImage(c, uint(request.ID), request.Size)
	if err != nil {
		log.Error().Err(err).Msgf("[card_controller:HandleImage] failed to get cached image")
	} else {
		c.Data(http.StatusOK, "image/jpeg", cached)
		return
	}

	img, err := ch.cardsRepository.GetImage(c, request.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dat, err := resizer.Resize(img.Data, request.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		if err := ch.cache.SaveImage(c, request.ID, request.Size, dat); err != nil {
			log.Error().Err(err).Send()
		}
	}()
	c.Data(http.StatusOK, "image/jpeg", dat)
}
