package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListRequest struct {
	From uint `uri:"from" binding:"required,numeric"`
	To   uint `uri:"to" binding:"required,numeric"`
}

type ListHandler struct {
	paintingsRepository paintingsRepository
}

func NewListHandler(paintingsRepository paintingsRepository) *ListHandler {
	return &ListHandler{paintingsRepository}
}

func (lh *ListHandler) Handle(c *gin.Context) {
	var request ListRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.From < 0 || request.To < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from or to must be positive numbers"})
	}
	paintings, err := lh.paintingsRepository.GetPaintingsRange(c, request.From, request.To)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, paintings)
}
