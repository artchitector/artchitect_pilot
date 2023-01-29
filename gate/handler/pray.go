package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
)

type PrayCreateRequest struct {
	Password string `json:"password" binding:"required"`
}

type PrayAnswerRequest struct {
	ID       uint   `json:"id" binding:"required,numeric"`
	Password string `json:"password" binding:"required"`
}

type PrayAnswer struct {
	Queue  uint
	State  string
	Answer uint
}

type PrayHandler struct {
	prayRepository prayRepository
}

func NewPrayHandler(prayRepository prayRepository) *PrayHandler {
	return &PrayHandler{prayRepository}
}

func (ph *PrayHandler) Handle(c *gin.Context) {
	r := PrayCreateRequest{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	log.Info().Msgf("[pray] incoming pray")
	pray, err := ph.prayRepository.MakePray(c, r.Password)
	if err != nil {
		log.Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "answers not available at the moment. sorry, maybe later."})
		return
	}
	log.Info().Msgf("[pray] saved pray %d", pray.ID)
	c.JSON(http.StatusOK, pray.ID)
}

func (ph *PrayHandler) HandleAnswer(c *gin.Context) {
	var request PrayAnswerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pray, err := ph.prayRepository.GetPrayWithPassword(c, request.ID, request.Password)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "pray not found or password incorrect"})
		log.Info().Msgf("[pray] failed to get pray %d with given password.", request.ID)
		return
	}
	if err != nil {
		log.Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "answers not available at the moment. sorry, maybe later."})
		return
	}

	queue, err := ph.prayRepository.GetQueueBeforePray(c, pray.ID)
	answer := PrayAnswer{
		Queue:  queue,
		State:  pray.State,
		Answer: pray.Answer,
	}
	c.JSON(http.StatusOK, answer)
}
