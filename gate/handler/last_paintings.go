package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LastPaintingsRequest struct {
	Quantity uint `uri:"quantity" binding:"required,numeric"`
}

type LastPaintingsHandler struct {
	paintingsRepository paintingsRepository
}

func NewLastPaintingsHandler(paintingsRepository paintingsRepository) *LastPaintingsHandler {
	return &LastPaintingsHandler{paintingsRepository: paintingsRepository}
}

func (lph *LastPaintingsHandler) Handle(c *gin.Context) {
	var request LastPaintingsRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quantity := request.Quantity
	if quantity > 100 {
		c.JSON(http.StatusBadGateway, "quantity required less than 100")
		return
	}
	paintings, err := lph.paintingsRepository.GetLastPaintings(c, uint64(quantity))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, paintings)
}
