package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type LikeRequest struct {
	CardID uint `uri:"card_id" json:"card_id" binding:"required"`
}

type LikedResponse struct {
	Liked bool
}

type bot interface {
	SendCardToInfinite(ctx context.Context, cardID uint, caption string) error
}

type LikeHandler struct {
	likeRepository       likeRepository
	cardsRepository      cardsRepository
	authService          *AuthService
	enhotter             enhotter
	artchitector         uint
	bot                  bot
	sendToInfiniteOnLike bool
}

func NewLikeHandler(
	likeRepository likeRepository,
	cardsRepository cardsRepository,
	authService *AuthService,
	enhotter enhotter,
	bot bot,
	artchitector uint,
	sendToInfiniteOnLike bool,
) *LikeHandler {
	return &LikeHandler{likeRepository, cardsRepository, authService, enhotter, artchitector, bot, sendToInfiniteOnLike}
}

func (lh *LikeHandler) Handle(c *gin.Context) {
	r := LikeRequest{}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	userID := lh.authService.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	like, err := lh.likeRepository.Like(c, userID, r.CardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Info().Msgf("userID is %d, artchitector is %d")
	go func(liked bool) {
		if liked {
			if err := lh.cardsRepository.Like(c, r.CardID); err != nil {
				log.Error().Err(err).Msgf("[like_handler] failed to like %d", r.CardID)
			}
		} else {
			if err := lh.cardsRepository.Unlike(c, r.CardID); err != nil {
				log.Error().Err(err).Msgf("[like_handler] failed to unlike %d", r.CardID)
			}
		}
		// update card cache
		log.Info().Msgf("GOROTINE#1")
		lh.enhotter.ReloadCardWithoutImage(c, r.CardID)
	}(like.Liked)
	if lh.sendToInfiniteOnLike && like.Liked && userID == lh.artchitector {
		// send this card to infinite
		go func() {
			if err := lh.bot.SendCardToInfinite(c, r.CardID, ""); err != nil {
				log.Error().Err(err).Msgf("[like_handler] failed send card %d to infite after like of %d", r.CardID, userID)
			} else {
				log.Info().Msgf("[like_handler] sent card %d to infinite after like of %d", r.CardID, userID)
			}
		}()
	}

	c.JSON(http.StatusOK, like)
}

func (lh *LikeHandler) HandleList(c *gin.Context) {
	userID := lh.authService.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	likes, err := lh.likeRepository.GetLikes(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, likes)
}

func (lh *LikeHandler) HandleGet(c *gin.Context) {
	var request LikeRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := lh.authService.getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	liked, err := lh.likeRepository.IsLiked(c, userID, request.CardID)
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lr := LikedResponse{Liked: liked}
	c.JSON(http.StatusOK, lr)
}
