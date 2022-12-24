package handler

import (
	"github.com/artchitector/artchitect.git/gate/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type StateResponse struct {
	Hash  string
	State model.State
}

type StateHandler struct {
	logger    zerolog.Logger
	retriever retriever
}

func NewStateHandler(logger zerolog.Logger, retriever retriever) *StateHandler {
	return &StateHandler{logger, retriever}
}

func (sh *StateHandler) Handle(c *gin.Context) {
	time.Sleep(time.Millisecond * 300) // slow down, like in 1990th
	state, err := sh.retriever.CollectState(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := StateResponse{Hash: "empty", State: state}
	c.JSON(http.StatusOK, response)
}
