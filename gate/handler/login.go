package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginHandler struct {
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (lh *LoginHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": ""})
}
