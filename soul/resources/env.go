package resources

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type Env struct {
	// enabled internal services
	LotteryEnabled      bool
	CardCreationEnabled bool
	MercifulEnabled     bool
	GifterActive        bool
	UseFakeArtist       bool
	TelegramBotEnabled  bool
	StorageEnabled      bool

	// external resources
	DbDSN          string
	RedisHostRU    string
	RedisHostEU    string
	RedisPassword  string
	OriginURL      string
	ArtistURL      string
	SaverURL       string
	MinioHost      string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string

	// settings
	CardTotalTime uint
	PrehotDelay   uint

	// telegram constants

	TelegramBotToken   string
	TenMinChat         string
	InfiniteChat       string
	ArtchitectorChatID int64
}

func initEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	cardTotalTimeStr := os.Getenv("CARD_TOTAL_TIME")
	cardTotalTime, err := strconv.Atoi(cardTotalTimeStr)
	if err != nil {
		log.Fatal().Err(err)
	}
	prehotDelayStr := os.Getenv("PREHOT_TIME")
	prehotDelay, err := strconv.Atoi(prehotDelayStr)
	if err != nil {
		log.Fatal().Err(err)
	}
	artchitectorChatStr := os.Getenv("ARTCHITECTOR_CHAT_ID")
	artchitectorChatID, err := strconv.ParseInt(artchitectorChatStr, 10, 64)

	if err != nil {
		log.Fatal().Err(err)
	}

	return &Env{
		LotteryEnabled:      os.Getenv("LOTTERY_ENABLED") == "true",
		CardCreationEnabled: os.Getenv("CARDS_CREATION_ENABLED") == "true",
		MercifulEnabled:     os.Getenv("MERCIFUL_ENABLED") == "true",
		GifterActive:        os.Getenv("GIFTER_ACTIVE") == "true",
		UseFakeArtist:       os.Getenv("USE_FAKE_ARTIST") == "true",
		TelegramBotEnabled:  os.Getenv("TELEGRAM_BOT_ENABLE") == "true",
		StorageEnabled:      os.Getenv("STORAGE_ENABLED") == "true",

		DbDSN:          os.Getenv("DB_DSN"),
		RedisHostRU:    os.Getenv("REDIS_HOST_RU"),
		RedisHostEU:    os.Getenv("REDIS_HOST_EU"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		OriginURL:      os.Getenv("ORIGIN_URL"),
		ArtistURL:      os.Getenv("ARTIST_URL"),
		MinioHost:      os.Getenv("MINIO_HOST"),
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinioBucket:    os.Getenv("MINIO_BUCKET"),
		SaverURL:       os.Getenv("SAVER_URL"),

		CardTotalTime: uint(cardTotalTime),
		PrehotDelay:   uint(prehotDelay),

		TelegramBotToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		TenMinChat:         os.Getenv("TEN_MIN_CHAT"),
		InfiniteChat:       os.Getenv("INFINITE_CHAT"),
		ArtchitectorChatID: artchitectorChatID,
	}
}
