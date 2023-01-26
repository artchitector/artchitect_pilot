package main

import (
	"context"
	artchitectService "github.com/artchitector/artchitect/soul/core/artchitect"
	artistService "github.com/artchitector/artchitect/soul/core/artist"
	engine2 "github.com/artchitector/artchitect/soul/core/artist/engine"
	creator2 "github.com/artchitector/artchitect/soul/core/creator"
	"github.com/artchitector/artchitect/soul/core/gifter"
	"github.com/artchitector/artchitect/soul/core/lottery"
	merciful2 "github.com/artchitector/artchitect/soul/core/merciful"
	originService "github.com/artchitector/artchitect/soul/core/origin"
	"github.com/artchitector/artchitect/soul/core/origin/driver"
	spellerService "github.com/artchitector/artchitect/soul/core/speller"
	notifier2 "github.com/artchitector/artchitect/soul/notifier"
	"github.com/artchitector/artchitect/soul/repository"
	"github.com/artchitector/artchitect/soul/resources"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	res := resources.InitResources()
	log.Info().Msg("[main] service soul started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	// repositoties
	cardsRepo := repository.NewCardRepository(res.GetDB())
	spellRepo := repository.NewSpellRepository(res.GetDB())
	lotteryRepo := repository.NewLotteryRepository(res.GetDB())
	prayRepo := repository.NewPrayRepository(res.GetDB())

	// notifier
	notifier := notifier2.NewNotifier(res.GetRedis())

	// origin
	webcamDriver := driver.NewWebcamDriver(res.GetEnv().OriginURL)
	origin := originService.NewOrigin(webcamDriver)

	// speller+artist+creator
	speller := spellerService.NewSpeller(spellRepo, origin, notifier)
	var engine artistService.EngineContract
	if res.GetEnv().UseFakeArtist {
		engine = engine2.NewFakeEngine()
	} else {
		engine = engine2.NewArtistEngine(res.GetEnv().ArtistURL)
	}
	artist := artistService.NewArtist(engine, cardsRepo, notifier)
	creator := creator2.NewCreator(artist, speller, notifier, res.GetEnv().CardTotalTime)

	// lottery runner
	runner := lottery.NewRunner(lotteryRepo, cardsRepo, origin)

	// merciful
	merciful := merciful2.NewMerciful(prayRepo, creator, notifier)

	// Artchitect core scheduler
	artchitectConfig := artchitectService.Config{
		CardsCreationEnabled: res.GetEnv().CardCreationEnabled,
		LotteryEnabled:       res.GetEnv().LotteryEnabled,
		MercifulEnabled:      res.GetEnv().MercifulEnabled,
	}
	artchitect := artchitectService.NewArtchitect(
		artchitectConfig,
		creator,
		lotteryRepo,
		runner,
		merciful,
		notifier,
	)

	// gifter
	if res.GetEnv().GifterActive {
		gift := gifter.NewGifter(cardsRepo, origin, res.GetEnv().TelegramBotToken, res.GetEnv().TenMinChat)
		go func() {
			if err := gift.Run(ctx); err != nil {
				log.Fatal().Err(err).Send()
			}
		}()
	}

	// main loop to make artworks
	var tick int
mainFor:
	for {
		select {
		case <-ctx.Done():
			break mainFor
		case <-time.Tick(time.Second * 1):
			tick += 1
			err := artchitect.Run(ctx, tick)
			if err != nil {
				log.Error().Err(err).Msgf("[main] failed to run artchitect task")
			}
		}
	}

	log.Info().Msg("[main] soul.Run finished")
}
