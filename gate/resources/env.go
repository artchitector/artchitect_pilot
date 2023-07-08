package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type Env struct {
	DbDSN          string
	HttpPort       string
	RedisHost      string
	RedisPassword  string
	MemoryHost     string
	JWTSecret      string
	ArtchitectHost string
	AllowFakeAuth  bool

	// telegram constants
	Telegram10BotToken   string // 10bot (is for maintenance and secure use to control artchitect.space). Secured with single account usage.
	TelegramABotToken    string // ABot (is for everyone: login, prayer etc)
	ChatID10             string
	ChatIDInfinite       string
	ChatIDArtchitector   int64
	SendToInfiniteOnLike bool
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	artchitectorChatStr := os.Getenv("CHAT_ID_ARTCHITECTOR")
	artchitectorChatID, err := strconv.ParseInt(artchitectorChatStr, 10, 64)

	return &Env{
		DbDSN:          os.Getenv("DB_DSN"),
		HttpPort:       os.Getenv("HTTP_PORT"),
		RedisHost:      os.Getenv("REDIS_HOST"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		MemoryHost:     os.Getenv("MEMORY_HOST"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		ArtchitectHost: os.Getenv("ARTCHITECT_HOST"),
		AllowFakeAuth:  os.Getenv("ALLOW_FAKE_AUTH") == "true",

		Telegram10BotToken:   os.Getenv("TELEGRAM_10BOT_TOKEN"),
		TelegramABotToken:    os.Getenv("TELEGRAM_ABOT_TOKEN"),
		ChatID10:             os.Getenv("CHAT_ID_10MIN"),
		ChatIDInfinite:       os.Getenv("CHAT_ID_INFINITE"),
		ChatIDArtchitector:   artchitectorChatID,
		SendToInfiniteOnLike: os.Getenv("SEND2INFINITE_ON_LIKE") == "true",
	}
}
