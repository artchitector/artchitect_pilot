package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type SelectionHandler struct {
	selectionRepository selectionRepository
}

func NewSelectionHandler(selectionRepository selectionRepository) *SelectionHandler {
	return &SelectionHandler{selectionRepository: selectionRepository}
}

func (lh *SelectionHandler) Handle(c *gin.Context) {
	selection, err := lh.selectionRepository.GetSelection(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().Msgf("%+v", selection)
	c.JSON(http.StatusOK, selection)
}
