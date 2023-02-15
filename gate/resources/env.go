package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type Env struct {
	DbDSN             string
	HttpPort          string
	RedisHost         string
	RedisPassword     string
	MemoryHost        string
	TelegramABotToken string
	JWTSecret         string
	ArtchitectHost    string
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return &Env{
		DbDSN:             os.Getenv("DB_DSN"),
		HttpPort:          os.Getenv("HTTP_PORT"),
		RedisHost:         os.Getenv("REDIS_HOST"),
		RedisPassword:     os.Getenv("REDIS_PASSWORD"),
		MemoryHost:        os.Getenv("MEMORY_HOST"),
		TelegramABotToken: os.Getenv("TELEGRAM_ABOT_TOKEN"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		ArtchitectHost:    os.Getenv("ARTCHITECT_HOST"),
	}
}
