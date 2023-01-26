package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type Env struct {
	// enabled internal services
	LotteryEnabled      bool
	CardCreationEnabled bool
	MercifulEnabled     bool
	GifterActive        bool
	UseFakeArtist       bool

	// external resources
	DbDSN         string
	RedisHost     string
	RedisPassword string
	OriginURL     string
	ArtistURL     string

	// telegram constants
	TelegramBotToken string
	TenMinChat       string
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return &Env{
		LotteryEnabled:      os.Getenv("LOTTERY_ENABLED") == "true",
		CardCreationEnabled: os.Getenv("CARDS_CREATION_ENABLED") == "true",
		MercifulEnabled:     os.Getenv("MERCIFUL_ENABLED") == "true",
		GifterActive:        os.Getenv("GIFTER_ACTIVE") == "true",
		UseFakeArtist:       os.Getenv("USE_FAKE_ARTIST") == "true",

		DbDSN:         os.Getenv("DB_DSN"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		OriginURL:     os.Getenv("ORIGIN_URL"),
		ArtistURL:     os.Getenv("ARTIST_URL"),

		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TenMinChat:       os.Getenv("TEN_MIN_CHAT"),
	}
}
