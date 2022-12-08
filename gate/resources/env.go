package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type Env struct {
	DbDSN    string
	HttpPort string
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return &Env{
		DbDSN:    os.Getenv("DB_DSN"),
		HttpPort: os.Getenv("HTTP_PORT"),
	}
}
