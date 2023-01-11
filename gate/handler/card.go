package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CardRequest struct {
	ID uint `uri:"id" binding:"required,numeric"`
}

type CardHandler struct {
	cardsRepository cardsRepository
}

func NewCardHandler(cardsRepository cardsRepository) *CardHandler {
	return &CardHandler{cardsRepository: cardsRepository}
}

func (lh *CardHandler) Handle(c *gin.Context) {
	var request CardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card, found, err := lh.cardsRepository.GetCard(c, request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, card)
}
