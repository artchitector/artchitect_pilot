package main

import (
	"context"
	cache2 "github.com/artchitector/artchitect/gate/cache"
	"github.com/artchitector/artchitect/gate/handler"
	"github.com/artchitector/artchitect/gate/listener"
	"github.com/artchitector/artchitect/gate/repository"
	"github.com/artchitector/artchitect/gate/resources"
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

	// cache
	cache := cache2.NewCache(res.GetRedis())

	// repos
	cardsRepo := repository.NewCardRepository(res.GetDB(), cache)
	lotteryRepo := repository.NewLotteryRepository(res.GetDB())
	prayRepo := repository.NewPrayRepository(res.GetDB())
	selectionRepo := repository.NewSelectionRepository(res.GetDB())

	// refresher (update cache)
	refresher := repository.NewRefresher(cardsRepo, selectionRepo)
	go func() {
		if err := refresher.RefreshLast(ctx); err != nil {
			log.Error().Err(err).Msgf("[main] failed refreshing last")
			cancel() // stop application and it will be reloaded
		}
		if err := refresher.RefreshSelection(ctx); err != nil {
			log.Error().Err(err).Msgf("[main] failed refreshing selection")
			cancel() // stop application and it will be reloaded
		}
	}()
	go func() {
		if err := refresher.StartRefreshing(ctx); err != nil {
			log.Error().Err(err).Msgf("[main] failed start refreshing")
			cancel() // stop application and it will be reloaded
		}
	}()

	// handlers
	lastPaintingsHandler := handler.NewLastPaintingsHandler(cardsRepo, cache)
	lotteryHandler := handler.NewLotteryHandler(
		log.With().Str("service", "lottery_handler").Logger(),
		lotteryRepo,
	)
	cardHandler := handler.NewCardHandler(cardsRepo, cache)
	selectionHander := handler.NewSelectionHandler(selectionRepo)
	prayHandler := handler.NewPrayHandler(prayRepo)

	// listeners with websocket handler
	lis := listener.NewListener(res.GetRedis(), cache, cardsRepo)
	websocketHandler := handler.NewWebsocketHandler(lis)

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
		r.GET("/last_paintings/:quantity", lastPaintingsHandler.Handle)
		r.GET("/lottery/:lastN", lotteryHandler.HandleLast)
		r.GET("/card/:id", cardHandler.Handle)
		r.GET("/selection", selectionHander.Handle)
		r.GET("/image/:size/:id", cardHandler.HandleImage)
		r.GET("/ws", func(c *gin.Context) {
			websocketHandler.Handle(c.Writer, c.Request)
		})
		r.POST("/pray", prayHandler.Handle)
		r.POST("/pray/answer", prayHandler.HandleAnswer)
		if err := r.Run("0.0.0.0:" + res.GetEnv().HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("gate.Run finished")
}
