package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type LoginHandler struct {
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (lh *LoginHandler) Handle(c *gin.Context) {
	log.Info().Msgf("query: %+v", c.Request.URL.Query())
	c.JSON(http.StatusOK, "ok")
}
