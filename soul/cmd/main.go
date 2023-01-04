package main

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/core/origin"
	"github.com/artchitector/artchitect.git/soul/core/origin/driver"
	stateService "github.com/artchitector/artchitect.git/soul/core/state"
	"github.com/artchitector/artchitect.git/soul/repository"
	"github.com/artchitector/artchitect.git/soul/resources"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})

	res := resources.InitResources()
	log.Info().Msg("service soul started")

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	//paintingRepo := repository.NewPaintingRepository(res.GetDB())
	decisionRepo := repository.NewDecisionRepository(res.GetDB())
	stateRepository := repository.NewStateRepository(res.GetDB())

	//randProvider := driver.NewRandDriver()
	webcamDriver := driver.NewWebcamDriver(res.GetEnv().OriginURL, decisionRepo)
	originInstance := origin.NewOrigin(webcamDriver)

	state := stateService.NewState(stateRepository)
	go func() {
		if err := state.Run(ctx); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			break
		case <-time.Tick(time.Second):
			yes, err := originInstance.YesNo(ctx)
			if err != nil {
				log.Error().Err(err).Send()
			} else {
				log.Info().Msgf("[main] origin answered yes=%t", yes)
			}
		}
	}

	log.Info().Msg("soul.Run finished")
}
