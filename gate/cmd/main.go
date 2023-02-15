package main

import (
	"context"
	cache2 "github.com/artchitector/artchitect/gate/cache"
	"github.com/artchitector/artchitect/gate/handler"
	"github.com/artchitector/artchitect/gate/listener"
	"github.com/artchitector/artchitect/gate/memory"
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

	// repos
	cardsRepo := repository.NewCardRepository(res.GetDB())
	lotteryRepo := repository.NewLotteryRepository(res.GetDB())
	prayRepo := repository.NewPrayRepository(res.GetDB())
	selectionRepo := repository.NewSelectionRepository(res.GetDB())

	// cache
	cache := cache2.NewCache(res.GetRedis())
	//_ = cache.Flushall(ctx)
	mmr := memory.NewMemory(res.GetEnv().MemoryHost, cache)
	enhotter := cache2.NewEnhotter(cardsRepo, selectionRepo, cache, mmr)
	enhotter.Run(ctx)

	// handlers
	lastCardsHandler := handler.NewLastCardsHandler(cardsRepo, cache)
	lotteryHandler := handler.NewLotteryHandler(
		log.With().Str("service", "lottery_handler").Logger(),
		lotteryRepo,
	)
	cardHandler := handler.NewCardHandler(cardsRepo, cache, mmr)
	selectionHander := handler.NewSelectionHandler(selectionRepo)
	prayHandler := handler.NewPrayHandler(prayRepo)
	lh := handler.NewLoginHandler(res.GetEnv().TelegramABotToken)

	// listeners with websocket handler
	lis := listener.NewListener(res.GetRedis(), cache, cardsRepo, mmr)
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
			AllowOrigins: []string{"http://localhost", "https://artchitect.space", "https://ru.artchitect.space", "https://eu.artchitect.space"},
		}))
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal().Err(err).Send()
		}
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		r.GET("/last_paintings/:quantity", lastCardsHandler.Handle)
		r.GET("/lottery/:lastN", lotteryHandler.HandleLast)
		r.GET("/card/:id", cardHandler.Handle)
		r.GET("/selection", selectionHander.Handle)
		r.GET("/image/:size/:id", cardHandler.HandleImage)
		r.GET("/ws", func(c *gin.Context) {
			websocketHandler.Handle(c.Writer, c.Request)
		})
		r.POST("/pray", prayHandler.Handle)
		r.POST("/pray/answer", prayHandler.HandleAnswer)
		r.GET("/login", lh.Handle)

		if err := r.Run("0.0.0.0:" + res.GetEnv().HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("gate.Run finished")
}
