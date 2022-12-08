package main

import (
	"context"
	"github.com/artchitector/artchitect.git/gate/resources"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_ = resources.InitResources()
	log.Info().Msg("service gate started")

	_, cancel := context.WithCancel(context.Background())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	log.Info().Msg("gate.Run finished")
}
