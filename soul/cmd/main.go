package main

import (
	"context"
	"github.com/artchitector/artchitect/model"
	artistService "github.com/artchitector/artchitect/soul/core/artist"
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

	paintingRepo := repository.NewPaintingRepository(res.GetDB())
	decisionRepo := repository.NewDecisionRepository(res.GetDB())
	stateRepository := repository.NewStateRepository(res.GetDB())
	spellRepository := repository.NewSpellRepository(res.GetDB())

	//randProvider := driver.NewRandDriver()
	webcamDriver := driver.NewWebcamDriver(res.GetEnv().OriginURL, decisionRepo)
	origin := originService.NewOrigin(webcamDriver)
	speller := spellerService.NewSpeller(spellRepository, origin)
	artist := artistService.NewArtist(res.GetEnv().ArtistURL, paintingRepo)

	state := stateService.NewState(stateRepository)

	if !res.GetEnv().Enabled {
		log.Info().Msg("[main] Soul service is not enabled")
		return
	}

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
		case <-time.Tick(time.Second * 10):
			state.SetState(ctx, model.StateMakingSpell)
			log.Info().Msgf("[main] start main loop")
			spell, err := speller.MakeSpell(ctx)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}
			log.Info().Msgf("[main] spell: %+v", spell)
			state.SetState(ctx, model.StateMakingArtifact)
			painting, err := artist.GetPainting(ctx, spell)
			if err != nil {
				log.Error().Err(err).Send()
				continue
			}
			log.Info().Msgf("[main] painting: id=%d, spell_id=%d", painting.ID, spell.ID)
			state.SetState(ctx, model.StateMakingRest)
		}
	}

	log.Info().Msg("soul.Run finished")
}
