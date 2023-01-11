package main

import (
	"context"
	artchitectService "github.com/artchitector/artchitect/soul/core/artchitect"
	artistService "github.com/artchitector/artchitect/soul/core/artist"
	"github.com/artchitector/artchitect/soul/core/lottery"
	originService "github.com/artchitector/artchitect/soul/core/origin"
	"github.com/artchitector/artchitect/soul/core/origin/driver"
	spellerService "github.com/artchitector/artchitect/soul/core/speller"
	stateService "github.com/artchitector/artchitect/soul/core/state"
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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	res := resources.InitResources()
	log.Info().Msg("service soul started")

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	paintingRepo := repository.NewCardRepository(res.GetDB())
	decisionRepo := repository.NewDecisionRepository(res.GetDB())
	stateRepository := repository.NewStateRepository(res.GetDB())
	spellRepository := repository.NewSpellRepository(res.GetDB())
	lotteryRepository := repository.NewLotteryRepository(res.GetDB())

	//randProvider := driver.NewRandDriver()
	webcamDriver := driver.NewWebcamDriver(res.GetEnv().OriginURL, decisionRepo)
	origin := originService.NewOrigin(webcamDriver)
	speller := spellerService.NewSpeller(spellRepository, origin)
	artist := artistService.NewArtist(res.GetEnv().ArtistURL, paintingRepo)
	runner := lottery.NewRunner(lotteryRepository, paintingRepo, origin)

	state := stateService.NewState(stateRepository)

	artchitectConfig := artchitectService.Config{
		CardsCreationEnabled: res.GetEnv().CardCreationEnabled,
		LotteryEnabled:       res.GetEnv().LotteryEnabled,
	}
	artchitect := artchitectService.NewArtchitect(artchitectConfig, state, speller, artist, lotteryRepository, runner)

	// state saving (in DB) process
	go func() {
		if err := state.Run(ctx); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	// main loop to make artworks
	for {
		select {
		case <-ctx.Done():
			break
		case <-time.Tick(time.Second * 1):
			err := artchitect.Run(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("failed to run artchitect task")
			}
		}
	}

	log.Info().Msg("soul.Run finished")
}
