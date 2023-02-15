package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LikeRequest struct {
	CardID uint `json:"card_id" binding:"required"`
}

type LikeHandler struct {
	likeRepository likeRepository
	authService    *AuthService
}

func NewLikeHandler(likeRepository likeRepository, authService *AuthService) *LikeHandler {
	return &LikeHandler{likeRepository, authService}
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
	c.JSON(http.StatusOK, like)
}
