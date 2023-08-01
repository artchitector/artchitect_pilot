package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type Env struct {
	// enabled internal services
	LotteryEnabled       bool
	CardCreationEnabled  bool
	MercifulEnabled      bool
	UnifierEnabled       bool
	GifterActive         bool
	UseFakeArtist        bool
	Telegram10BotEnabled bool
	TelegramABotEnabled  bool
	StorageEnabled       bool

	// external resources
	DbDSN           string
	RedisHostRU     string
	RedisHostEU     string
	RedisPassword   string
	OriginURL       string
	ArtistURL       string
	MemorySaverURL  string
	MemoryHost      string
	StorageSaverURL string

	// settings
	ArtTotalTime       uint
	PrehotDelay        uint
	FakeGenerationTime uint

	// telegram constants
	Telegram10BotToken string // 10bot (is for maintenance and secure use to control artchitect.space). Secured with single account usage.
	TelegramABotToken  string // ABot (is for everyone: login, prayer etc)
	ChatID10           string
	ChatIDInfinite     string
	ChatIDArtchitector int64
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	artTotalTimeStr := os.Getenv("ART_TOTAL_TIME")
	artTotalTime, err := strconv.Atoi(artTotalTimeStr)
	if err != nil {
		log.Fatal().Err(err)
	}
	prehotDelayStr := os.Getenv("PREHOT_TIME")
	prehotDelay, err := strconv.Atoi(prehotDelayStr)
	if err != nil {
		log.Fatal().Err(err)
	}
	artchitectorChatStr := os.Getenv("CHAT_ID_ARTCHITECTOR")
	artchitectorChatID, err := strconv.ParseInt(artchitectorChatStr, 10, 64)

	fakeGenerationTimeStr := os.Getenv("FAKE_GENERATION_TIME")
	fakeGenerationTime, err := strconv.Atoi(fakeGenerationTimeStr)
	if err != nil {
		log.Fatal().Err(err)
	}

	if err != nil {
		log.Fatal().Err(err)
	}

	return &Env{
		LotteryEnabled:       os.Getenv("LOTTERY_ENABLED") == "true",
		CardCreationEnabled:  os.Getenv("CARDS_CREATION_ENABLED") == "true",
		MercifulEnabled:      os.Getenv("MERCIFUL_ENABLED") == "true",
		UnifierEnabled:       os.Getenv("UNIFIER_ENABLED") == "true",
		GifterActive:         os.Getenv("GIFTER_ACTIVE") == "true",
		UseFakeArtist:        os.Getenv("USE_FAKE_ARTIST") == "true",
		Telegram10BotEnabled: os.Getenv("TELEGRAM_10BOT_ENABLE") == "true",
		TelegramABotEnabled:  os.Getenv("TELEGRAM_ABOT_ENABLE") == "true",
		StorageEnabled:       os.Getenv("STORAGE_ENABLED") == "true",

		DbDSN:           os.Getenv("DB_DSN"),
		RedisHostRU:     os.Getenv("REDIS_HOST_RU"),
		RedisHostEU:     os.Getenv("REDIS_HOST_EU"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		OriginURL:       os.Getenv("ORIGIN_URL"),
		ArtistURL:       os.Getenv("ARTIST_URL"),
		MemoryHost:      os.Getenv("MEMORY_HOST"),
		MemorySaverURL:  os.Getenv("MEMORY_SAVER_URL"),
		StorageSaverURL: os.Getenv("STORAGE_SAVER_URL"),

		ArtTotalTime:       uint(artTotalTime),
		PrehotDelay:        uint(prehotDelay),
		FakeGenerationTime: uint(fakeGenerationTime),

		Telegram10BotToken: os.Getenv("TELEGRAM_10BOT_TOKEN"),
		TelegramABotToken:  os.Getenv("TELEGRAM_ABOT_TOKEN"),
		ChatID10:           os.Getenv("CHAT_ID_10MIN"),
		ChatIDInfinite:     os.Getenv("CHAT_ID_INFINITE"),
		ChatIDArtchitector: artchitectorChatID,
	}
}
