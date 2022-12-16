package main

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/core/artchitector"
	"github.com/artchitector/artchitect.git/soul/core/artist"
	"github.com/artchitector/artchitect.git/soul/infrastructure"
	"github.com/artchitector/artchitect.git/soul/repository"
	"github.com/artchitector/artchitect.git/soul/resources"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
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

	paintingRepo := repository.NewPaintingRepository(res.GetDB())

	cloud := infrastructure.NewCloud(log.With().Str("service", "cloud").Logger())
	art := artist.NewArtist(
		log.With().Str("service", "artist").Logger(),
		cloud,
	)
	if err := art.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("artist.Run failed")
	}
	schedule := artchitector.NewSchedule(log.With().Str("service", "schedule").Logger())
	artchitect := artchitector.NewArtchitect(
		log.With().Str("service", "artchitector").Logger(),
		schedule,
		cloud,
		paintingRepo,
	)

	if err := artchitect.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("artchitect.Run failed")
	}

	log.Info().Msg("soul.Run finished")
}
