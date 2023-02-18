package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LikeRequest struct {
	CardID uint `json:"card_id" binding:"required"`
}

type bot interface {
	SendCardToInfinite(ctx context.Context, cardID uint, caption string) error
}

type LikeHandler struct {
	likeRepository likeRepository
	authService    *AuthService
	artchitector   uint
	bot            bot
}

func NewLikeHandler(likeRepository likeRepository, authService *AuthService, bot bot, artchitector uint) *LikeHandler {
	return &LikeHandler{likeRepository, authService, artchitector, bot}
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
	if like.Liked && userID == lh.artchitector {
		// send this card to infinite
		if err := lh.bot.SendCardToInfinite(c, r.CardID, ""); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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
