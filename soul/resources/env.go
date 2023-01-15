package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

type Env struct {
	LotteryEnabled      bool
	CardCreationEnabled bool
	MercifulEnabled     bool
	DbDSN               string
	OriginURL           string
	ArtistURL           string
	TelegramBotToken    string
	GifterActive        bool
	TenMinChat          string
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	lotteryEnabledFlag := os.Getenv("LOTTERY_ENABLED")
	cardsEnabledFlag := os.Getenv("CARDS_CREATION_ENABLED")
	return &Env{
		LotteryEnabled:      lotteryEnabledFlag == "true" || lotteryEnabledFlag == "TRUE",
		CardCreationEnabled: cardsEnabledFlag == "true" || cardsEnabledFlag == "TRUE",
		DbDSN:               os.Getenv("DB_DSN"),
		OriginURL:           os.Getenv("ORIGIN_URL"),
		ArtistURL:           os.Getenv("ARTIST_URL"),
		TelegramBotToken:    os.Getenv("TELEGRAM_BOT_TOKEN"),
		GifterActive:        os.Getenv("GIFTER_ACTIVE") == "true",
		TenMinChat:          os.Getenv("10MIN_CHAT"),
		MercifulEnabled:     os.Getenv("MERCIFUL_ENABLED") == "true",
	}
}
