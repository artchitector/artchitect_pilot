package handler

import (
	"encoding/json"
	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type LotteryHandler struct {
	logger            zerolog.Logger
	lotteryRepository lotteryRepository
}

type LotteryRequest struct {
	LastN uint `uri:"lastN" binding:"required,numeric"`
}

func NewLotteryHandler(logger zerolog.Logger, lotteryRepository lotteryRepository) *LotteryHandler {
	return &LotteryHandler{logger, lotteryRepository}
}

func (lh *LotteryHandler) HandleLast(c *gin.Context) {
	var request LotteryRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if request.LastN > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 10 lotteries"})
		return
	}
	lotteries, err := lh.lotteryRepository.GetLastLotteries(c, request.LastN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for lotteryIdx, lottery := range lotteries {
		if lottery.State != model.LotteryStateRunning && lottery.State != model.LotteryStateFinished {
			continue
		}

		winners := []uint{}
		if err := json.Unmarshal([]byte(lottery.WinnersJSON), &winners); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		lotteries[lotteryIdx].Winners = winners
	}

	c.JSON(http.StatusOK, lotteries)
}
