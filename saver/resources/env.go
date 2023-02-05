package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type Env struct {
	HttpPort  string
	DbDSN     string
	CardsPath string
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return &Env{
		HttpPort:  os.Getenv("HTTP_PORT"),
		DbDSN:     os.Getenv("DB_DSN"),
		CardsPath: os.Getenv("CARDS_PATH"),
	}
}
