package main

import (
	"context"
	"github.com/artchitector/artchitect.git/gate/handler"
	"github.com/artchitector/artchitect.git/gate/repository"
	"github.com/artchitector/artchitect.git/gate/resources"
	"github.com/artchitector/artchitect.git/gate/state"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	res := resources.InitResources()
	log.Info().Msg("service gate started")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	paintingRepository := repository.NewPaintingRepository(res.GetDB())
	retriever := state.NewRetriever(
		log.With().Str("service", "retriever").Logger(),
		paintingRepository,
	)
	stateHandler := handler.NewStateHandler(
		log.With().Str("service", "state_handler").Logger(),
		retriever,
	)
	paintingHandler := handler.NewPaintingHandler(
		log.With().Str("service", "painting_handler").Logger(),
		retriever,
	)

	go func() {
		r := gin.Default()
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal().Err(err).Send()
		}
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		r.GET("/state", stateHandler.Handle)
		r.GET("/painting/:id", paintingHandler.Handle)

		if err := r.Run("0.0.0.0:" + res.GetEnv().HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("gate.Run finished")
}
