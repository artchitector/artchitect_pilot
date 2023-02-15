package main

import (
	"context"
	"github.com/artchitector/artchitect/community/bot"
	"github.com/artchitector/artchitect/community/resources"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	res := resources.InitResources()
	log.Info().Msg("[main] service soul started")

	b := bot.NewBot(res.GetEnv().BotToken)
	if err := b.Run(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}
}
