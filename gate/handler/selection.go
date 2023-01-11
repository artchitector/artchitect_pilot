package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SelectionHandler struct {
	lotteryRepository lotteryRepository
}

func NewSelectionHandler(lotteryRepository lotteryRepository) *SelectionHandler {
	return &SelectionHandler{lotteryRepository: lotteryRepository}
}

func (lh *SelectionHandler) Handle(c *gin.Context) {
	selection, err := lh.lotteryRepository.GetSelection(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, selection)
}
