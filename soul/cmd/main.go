package main

import (
	"context"
	"github.com/artchitector/artchitect.git/soul/core/artchitector"
	"github.com/artchitector/artchitect.git/soul/core/artist"
	"github.com/artchitector/artchitect.git/soul/infrastructure"
	"github.com/artchitector/artchitect.git/soul/resources"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	resources.InitResources()
	log.Info().Msg("service started")

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	//go func() {
	//	for {
	//		payload := struct {
	//			Payload int64 `json:"payload"`
	//		}{
	//			Payload: int64(rand.Intn(1000)),
	//		}
	//		body, _ := json.Marshal(payload)
	//		res.GetDB().Exec(fmt.Sprintf("NOTIFY events, '%s';", string(body)))
	//		time.Sleep(time.Second)
	//	}
	//}()

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
	)

	if err := artchitect.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("artchitect.Run failed")
	}

	log.Info().Msg("artchitect.Run finished. System shutdown")
}
