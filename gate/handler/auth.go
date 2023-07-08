package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

type AuthService struct {
	secretKey     []byte
	allowFakeAuth bool
}

func NewAuthService(secretKey string, allowFakeAuth bool) *AuthService {
	return &AuthService{[]byte(secretKey), allowFakeAuth}
}

func (as *AuthService) getUserID(c *gin.Context) uint {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "FAKE_LOCAL_TOKEN" {
		if !as.allowFakeAuth {
			log.Error().Msgf("[auth] Authorization token is FAKE_LOCAL_TOKEN, but fake is disabled")
			return 0
		}
		log.Info().Msgf("[auth] used 999 user_id for fake auth")
		return 999
	}

	if authHeader == "" {
		return 0
	}
	log.Info().Msgf("header: %s", authHeader)
	token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("[auth] unexpected singing method: %v", token.Header["alg"])
		}
		return as.secretKey, nil
	})
	if err != nil {
		log.Error().Err(err).Msgf("[auth] failed to parse token")
		return 0
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Info().Msgf("[auth] %+v", claims)
		ID, err := strconv.ParseInt(fmt.Sprintf("%s", claims["id"]), 10, 64)
		if err != nil {
			log.Error().Err(err).Msgf("[auth]")
			return 0
		}
		log.Info().Msgf("[auth] logged in %d", ID)
		return uint(ID)
	} else {
		log.Error().Msgf("[auth] not ok!")
		return 0
	}

	return 0
}
