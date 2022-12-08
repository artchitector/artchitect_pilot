package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type PaintingRequest struct {
	ID uint `uri:"id" binding:"required,numeric"`
}

type PaintingHandler struct {
	logger    zerolog.Logger
	retriever retriever
}

func NewPaintingHandler(logger zerolog.Logger, retriever retriever) *PaintingHandler {
	return &PaintingHandler{logger, retriever}
}

func (ph *PaintingHandler) Handle(c *gin.Context) {
	var request PaintingRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	bytes, err := ph.retriever.GetPaintingData(c, request.ID)
	if err != nil {
		// TODO make 404 response, if no painting
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Data(http.StatusOK, "image/jpeg", bytes)
}
