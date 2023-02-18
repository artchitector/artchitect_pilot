package main

import (
	"context"
	"github.com/artchitector/artchitect/bot"
	cache2 "github.com/artchitector/artchitect/gate/cache"
	"github.com/artchitector/artchitect/gate/fake"
	"github.com/artchitector/artchitect/gate/handler"
	"github.com/artchitector/artchitect/gate/listener"
	"github.com/artchitector/artchitect/gate/resources"
	"github.com/artchitector/artchitect/memory"
	repository2 "github.com/artchitector/artchitect/model/repository"
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
	cardsRepo := repository2.NewCardRepository(res.GetDB(), &fake.FakeOrigin{})
	lotteryRepo := repository2.NewLotteryRepository(res.GetDB())
	prayRepo := repository2.NewPrayRepository(res.GetDB())
	selectionRepo := repository2.NewSelectionRepository(res.GetDB())
	likeRepo := repository2.NewLikeRepository(res.GetDB())
	unityRepo := repository2.NewUnityRepository(res.GetDB())

	// cache
	cache := cache2.NewCache(res.GetRedis())
	//_ = cache.Flushall(ctx)
	mmr := memory.NewMemory(res.GetEnv().MemoryHost, cache)
	enhotter := cache2.NewEnhotter(cardsRepo, selectionRepo, cache, mmr)
	enhotter.Run(ctx)

	artchitectBot := bot.NewBot(
		res.GetEnv().Telegram10BotToken,
		cardsRepo,
		mmr,
		res.GetEnv().ChatIDArtchitector,
		res.GetEnv().ChatID10,
		res.GetEnv().ChatIDInfinite,
	)

	// handlers
	lastCardsHandler := handler.NewLastCardsHandler(cardsRepo, cache)
	lotteryHandler := handler.NewLotteryHandler(
		log.With().Str("service", "lottery_handler").Logger(),
		lotteryRepo,
	)
	authS := handler.NewAuthService(res.GetEnv().JWTSecret)
	cardHandler := handler.NewCardHandler(cardsRepo, cache, likeRepo, authS)
	selectionHander := handler.NewSelectionHandler(selectionRepo)
	prayHandler := handler.NewPrayHandler(prayRepo)
	lh := handler.NewLoginHandler(res.GetEnv().TelegramABotToken, res.GetEnv().JWTSecret, res.GetEnv().ArtchitectHost)
	llh := handler.NewLikeHandler(likeRepo, authS, artchitectBot, uint(res.GetEnv().ChatIDArtchitector))
	uh := handler.NewUnityHandler(unityRepo, cardsRepo)
	ih := handler.NewImageHandler(mmr)

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
		r.GET("/image/:size/:id", ih.HandleImage)
		r.GET("/image/unity/:mask/:size", ih.HandleUnity)
		r.GET("/ws", func(c *gin.Context) {
			websocketHandler.Handle(c.Writer, c.Request)
		})
		r.POST("/pray", prayHandler.Handle)
		r.POST("/pray/answer", prayHandler.HandleAnswer)
		r.GET("/login", lh.Handle)
		r.POST("/like", llh.Handle)
		r.GET("/liked", llh.HandleList)
		r.GET("/unity", uh.HandleList)
		r.GET("/unity/:mask", uh.HandleUnity)

		if err := r.Run("0.0.0.0:" + res.GetEnv().HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("gate.Setup finished")
}
