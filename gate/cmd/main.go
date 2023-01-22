package main

import (
	"context"
	cache2 "github.com/artchitector/artchitect/gate/cache"
	"github.com/artchitector/artchitect/gate/handler"
	"github.com/artchitector/artchitect/gate/listener"
	"github.com/artchitector/artchitect/gate/repository"
	"github.com/artchitector/artchitect/gate/resources"
	"github.com/artchitector/artchitect/gate/state"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})

	res := resources.InitResources()
	log.Info().Msg("service gate started")

	cache := cache2.NewCache(res.GetRedis())
	if err := cache.Flushall(ctx); err != nil {
		log.Error().Err(err).Msgf("[main] failed redis flushall")
	}

	cardsRepository := repository.NewCardRepository(res.GetDB(), cache)
	decisionRepository := repository.NewDecisionRepository(res.GetDB())
	stateRepository := repository.NewStateRepository(res.GetDB())
	spellRepository := repository.NewSpellRepository(res.GetDB())
	lotteryRepository := repository.NewLotteryRepository(res.GetDB())
	retriever := state.NewRetriever(
		log.With().Str("service", "retriever").Logger(),
		cardsRepository,
		decisionRepository,
		stateRepository,
		spellRepository,
	)

	// refresher enhot cache
	refresher := repository.NewRefresher(cardsRepository, lotteryRepository)
	go func() {
		if err := refresher.RefreshLast(ctx); err != nil {
			log.Error().Err(err).Msgf("[main] failed refreshing")
			cancel() // stop application and it will be reloaded
		}
	}()
	go func() {
		if err := refresher.StartRefreshing(ctx); err != nil {
			log.Error().Err(err).Msgf("[main] failed refreshing")
			cancel() // stop application and it will be reloaded
		}
	}()

	stateHandler := handler.NewStateHandler(
		log.With().Str("service", "state_handler").Logger(),
		retriever,
	)
	lastPaintingsHandler := handler.NewLastPaintingsHandler(cardsRepository, cache)
	lotteryHandler := handler.NewLotteryHandler(
		log.With().Str("service", "lottery_handler").Logger(),
		lotteryRepository,
	)
	cardHandler := handler.NewCardHandler(cardsRepository, cache)
	selectionHander := handler.NewSelectionHandler(lotteryRepository)

	prayRepository := repository.NewPrayRepository(res.GetDB())
	prayHandler := handler.NewPrayHandler(prayRepository)
	lis := listener.NewListener(res.GetRedis(), cache, cardsRepository)

	go func() {
		err := lis.Run(ctx)
		if err != nil {
			log.Error().Err(err).Send()
			cancel()
		}
	}()

	go func() {
		r := gin.Default()
		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
			//AllowOrigins:           nil,
		}))
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal().Err(err).Send()
		}
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		r.GET("/state", stateHandler.Handle)
		r.GET("/last_paintings/:quantity", lastPaintingsHandler.Handle)
		r.GET("/lottery/:lastN", lotteryHandler.HandleLast)
		r.GET("/card/:id", cardHandler.Handle)
		r.GET("/selection", selectionHander.Handle)
		r.GET("/image/:size/:id", cardHandler.HandleImage)
		r.GET("/answer", prayHandler.Handle)
		r.GET("/answer/:id", prayHandler.HandleAnswer)

		if err := r.Run("0.0.0.0:" + res.GetEnv().HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("gate.Run finished")
}
