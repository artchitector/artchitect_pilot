package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type LastCardsRequest struct {
	Quantity uint `uri:"quantity" binding:"required,numeric"`
}

type LastCardsHandler struct {
	artsRepository artsRepository
	cache          cache
}

func NewLastCardsHandler(paintingsRepository artsRepository, cache cache) *LastCardsHandler {
	return &LastCardsHandler{paintingsRepository, cache}
}

func (lph *LastCardsHandler) Handle(c *gin.Context) {
	var request LastCardsRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quantity := request.Quantity
	if quantity > 100 {
		c.JSON(http.StatusBadGateway, "quantity required less than 100")
		return
	}
	cards, err := lph.cache.GetLastCards(c, uint(quantity))
	if err != nil {
		log.Error().Err(err).Msgf("[last_cards_handler] failed to get last cards from cache")
	} else {
		c.JSON(http.StatusOK, cards)
		return
	}
	cards, err = lph.artsRepository.GetLastArts(c, uint(quantity))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, cards)
}
