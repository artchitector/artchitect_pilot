package main

import (
	"context"
	artchitectService "github.com/artchitector/artchitect/soul/core/artchitect"
	artistService "github.com/artchitector/artchitect/soul/core/artist"
	"github.com/artchitector/artchitect/soul/core/gifter"
	"github.com/artchitector/artchitect/soul/core/lottery"
	merciful2 "github.com/artchitector/artchitect/soul/core/merciful"
	originService "github.com/artchitector/artchitect/soul/core/origin"
	"github.com/artchitector/artchitect/soul/core/origin/driver"
	spellerService "github.com/artchitector/artchitect/soul/core/speller"
	stateService "github.com/artchitector/artchitect/soul/core/state"
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
	log.Info().Msg("service soul started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	cardsRepo := repository.NewCardRepository(res.GetDB())
	decisionRepo := repository.NewDecisionRepository(res.GetDB())
	stateRepository := repository.NewStateRepository(res.GetDB())
	spellRepository := repository.NewSpellRepository(res.GetDB())
	lotteryRepository := repository.NewLotteryRepository(res.GetDB())
	prayRepository := repository.NewPrayRepository(res.GetDB())

	//randProvider := driver.NewRandDriver()
	webcamDriver := driver.NewWebcamDriver(res.GetEnv().OriginURL, decisionRepo)
	origin := originService.NewOrigin(webcamDriver)
	speller := spellerService.NewSpeller(spellRepository, origin)
	var artist artchitectService.ArtistContract
	if res.GetEnv().UseFakeArtist {
		artist = artistService.NewFakeArtist(cardsRepo)
	} else {
		artist = artistService.NewArtist(res.GetEnv().ArtistURL, cardsRepo)
	}
	runner := lottery.NewRunner(lotteryRepository, cardsRepo, origin)
	state := stateService.NewState(stateRepository)
	merciful := merciful2.NewMerciful(prayRepository, artist, state, speller)

	notifier := notifier2.NewNotifier(res.GetRedis())

	artchitectConfig := artchitectService.Config{
		CardsCreationEnabled: res.GetEnv().CardCreationEnabled,
		LotteryEnabled:       res.GetEnv().LotteryEnabled,
		MercifulEnabled:      res.GetEnv().MercifulEnabled,
	}
	artchitect := artchitectService.NewArtchitect(
		artchitectConfig,
		state,
		speller,
		artist,
		lotteryRepository,
		runner,
		merciful,
		notifier,
	)
	gift := gifter.NewGifter(cardsRepo, origin, res.GetEnv().TelegramBotToken, res.GetEnv().TenMinChat)

	// state saving (in DB) process
	go func() {
		if err := state.Run(ctx); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	if res.GetEnv().GifterActive {
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
				log.Error().Err(err).Msgf("failed to run artchitect task")
			}
		}
	}

	log.Info().Msg("soul.Run finished")
}
