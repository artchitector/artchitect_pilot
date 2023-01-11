package handler

import (
	"encoding/json"
	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
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

func (lh *LotteryHandler) Handle(c *gin.Context) {
	lottery, err := lh.lotteryRepository.GetActiveLottery(c)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active lottery"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lottery)
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
		for tourIdx, tour := range lottery.Tours {
			winners := []uint64{}
			if err := json.Unmarshal([]byte(tour.WinnersJSON), &winners); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			lotteries[lotteryIdx].Tours[tourIdx].Winners = winners
			if lottery.ID == 4 && tour.ID == 8 {
				log.Info().Msgf("%+v", lotteries[lotteryIdx].Tours[tourIdx].Winners)
			}
		}

		if lottery.State != model.LotteryStateFinished {
			continue
		}

		winners := []uint64{}
		if err := json.Unmarshal([]byte(lottery.WinnersJSON), &winners); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		lotteries[lotteryIdx].Winners = winners
	}

	log.Info().Msgf("%+v", lotteries[1].Tours[0])

	c.JSON(http.StatusOK, lotteries)
}
