package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type Env struct {
	Enabled   bool
	DbDSN     string
	OriginURL string
	ArtistURL string
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	enabledFlag := os.Getenv("ENABLED")
	return &Env{
		Enabled:   enabledFlag == "true" || enabledFlag == "TRUE",
		DbDSN:     os.Getenv("DB_DSN"),
		OriginURL: os.Getenv("ORIGIN_URL"),
		ArtistURL: os.Getenv("ARTIST_URL"),
	}
}
