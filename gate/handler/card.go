package handler

import (
	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
)

type CardRequest struct {
	ID uint `uri:"id" binding:"required,numeric"`
}

type CardHandler struct {
	cardsRepository cardsRepository
	cache           cache
	likeRepository  likeRepository
	authService     *AuthService
}

func NewCardHandler(cardsRepository cardsRepository, cache cache, likeRepository likeRepository, authService *AuthService) *CardHandler {
	return &CardHandler{cardsRepository, cache, likeRepository, authService}
}

func (ch *CardHandler) Handle(c *gin.Context) {
	var request CardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var card model.Card
	var err error
	card, err = ch.cache.GetCard(c, uint(request.ID))
	if err != nil {
		log.Error().Err(err).Msgf("[card_handler:Handle] failed to get card(id=%d) from cache", request.ID)
	}

	if card.ID == 0 {
		card, err = ch.cardsRepository.GetCard(c, request.ID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	userID := ch.authService.getUserID(c)
	log.Info().Msgf("%d", userID)
	if userID != 0 {
		liked, err := ch.likeRepository.IsLiked(c, userID, card.ID)
		if err != nil {
			log.Error().Err(err).Msgf("[card_handler:Handle] failed to get liked")
		} else {
			card.Liked = liked
		}
	}

	c.JSON(http.StatusOK, card)
}
