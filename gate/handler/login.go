package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type LoginHandler struct {
	botToken       string
	secretKey      []byte
	artchitectHost string
}

func NewLoginHandler(botToken string, secretKey string, artchitectHost string) *LoginHandler {
	return &LoginHandler{botToken, []byte(secretKey), artchitectHost}
}

func (lh *LoginHandler) Handle(c *gin.Context) {
	log.Info().Msgf("query: %+v", c.Request.URL.Query())
	if err := lh.checkFromTelegram(c.Request.URL.Query()); err != nil {
		log.Error().Err(err).Msgf("[login_handler] failed telegram check")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed telegram check. you're not from telegram"})
		return
	}

	j, err := lh.generateJWT(c.Request.URL.Query())
	if err != nil {
		log.Error().Err(err).Msgf("[login_handler] failed to generate jwt")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt not generated"})
		return
	}

	params := url.Values{}
	params.Add("token", j)
	params.Add("username", c.Request.URL.Query().Get("username"))
	params.Add("photo_url", c.Request.URL.Query().Get("photo_url"))
	c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?%s", lh.artchitectHost, params.Encode()))
}

func (lh *LoginHandler) checkFromTelegram(values url.Values) error {
	/*
		https://core.telegram.org/widgets/login
		Data-check-string is a concatenation of all received fields, sorted in alphabetical order, in the format
		key=<value> with a line feed character ('\n', 0x0A) used as separator – e.g.,
		'auth_date=<auth_date>\nfirst_name=<first_name>\nid=<id>\nusername=<username>'.

		data_check_string = ...
		secret_key = SHA256(<bot_token>)
		if (hex(HMAC_SHA256(data_check_string, secret_key)) == hash) {
		  // data is from Telegram
		}
	*/
	hash := values.Get("hash")
	values.Del("hash")

	keys := make([]string, 0, len(values))
	for key, _ := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	stringPieces := make([]string, 0, len(values))

	for _, key := range keys {
		stringPieces = append(stringPieces, fmt.Sprintf("%s=%s", key, values.Get(key)))
	}
	dataCheckString := strings.Join(stringPieces, "\n")
	secretKey := makeSha256(lh.botToken)
	encryptedDataCheckString := makeHmacSha256([]byte(dataCheckString), secretKey)
	if hash != hex.EncodeToString(encryptedDataCheckString) {
		return errors.Errorf("[login_handler] hash not valid")
	}
	return nil
}

func makeSha256(str string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return hasher.Sum(nil)
}

func makeHmacSha256(data []byte, key []byte) []byte {
	sig := hmac.New(sha256.New, key)
	sig.Write(data)
	return sig.Sum(nil)
}

func (lh *LoginHandler) generateJWT(v url.Values) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = v.Get("id")
	claims["first_name"] = v.Get("first_name")
	claims["username"] = v.Get("username")
	claims["photo_url"] = v.Get("photo_url")
	claims["auth_date"] = v.Get("auth_date")

	log.Info().Msgf("%s", lh.secretKey)
	tokenStr, err := token.SignedString(lh.secretKey)
	if err != nil {
		return "", errors.Wrapf(err, "[login_handler] failed to sign JWT")
	}
	return tokenStr, nil
}
