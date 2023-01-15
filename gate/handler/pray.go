package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type PrayHandler struct {
	prayRepository prayRepository
}

func NewPrayHandler(prayRepository prayRepository) *PrayHandler {
	return &PrayHandler{prayRepository}
}

func (ph *PrayHandler) Handle(c *gin.Context) {
	log.Info().Msgf("[pray] incoming pray")
	pray, err := ph.prayRepository.MakePray(c)
	if err != nil {
		log.Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "answers not available at the moment. sorry, maybe later."})
		return
	}
	log.Info().Msgf("[pray] saved pray %d", pray.ID)
	time.Sleep(time.Second * 5)
	answer, err := ph.prayRepository.GetAnswer(c, uint64(pray.ID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, "0")
		log.Info().Msgf("[pray] failed to wait pray %d. 0 sent", pray.ID)
		return
	}
	if err != nil {
		log.Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "answers not available at the moment. sorry, maybe later."})
		return
	}
	c.JSON(http.StatusOK, answer)
	log.Info().Msgf("[pray] answered pray %d", pray.ID)
}
