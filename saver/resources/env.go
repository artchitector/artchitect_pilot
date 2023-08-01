package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type Env struct {
	HttpPort     string
	DbDSN        string
	ArtsPath     string
	UnityPath    string
	FullSizePath string
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return &Env{
		HttpPort:     os.Getenv("HTTP_PORT"),
		DbDSN:        os.Getenv("DB_DSN"),
		ArtsPath:     os.Getenv("ARTS_PATH"),
		UnityPath:    os.Getenv("UNITY_PATH"),
		FullSizePath: os.Getenv("FULLSIZE_PATH"),
	}
}
