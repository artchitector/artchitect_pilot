package handler

import (
	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HundredRequest struct {
	Hundred uint `uri:"hundred" binding:"required,numeric"`
}

type hundredRepository interface {
	FindAllTenK() ([]model.Hundred, error)
	FindKList(tenKHundred uint) ([]model.Hundred, error)
	FindHList(kHundred uint) ([]model.Hundred, error)
}
type cardRepository interface {
	GetHundred(hundred uint) ([]model.Card, error)
}
type SearchHandler struct {
	hundredRepository hundredRepository
	cardRepository    cardRepository
}

func NewSearchHandler(hundredRepository hundredRepository, cardRepository cardRepository) *SearchHandler {
	return &SearchHandler{hundredRepository, cardRepository}
}

// no params - show all 10k cards
func (sh *SearchHandler) HandleTenKList(c *gin.Context) {
	hundreds, err := sh.hundredRepository.FindAllTenK()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hundreds)
}

// params - show 1k cards in some 10k
func (sh *SearchHandler) HandleKList(c *gin.Context) {
	var request HundredRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hundreds, err := sh.hundredRepository.FindKList(request.Hundred)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hundreds)
}

// params - show 100 cards in some 1k
func (sh *SearchHandler) HandleHList(c *gin.Context) {
	var request HundredRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hundreds, err := sh.hundredRepository.FindHList(request.Hundred)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hundreds)
}

// params - show 100 cards in some 1k
func (sh *SearchHandler) HandleH(c *gin.Context) {
	var request HundredRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hundreds, err := sh.cardRepository.GetHundred(request.Hundred)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hundreds)
}
