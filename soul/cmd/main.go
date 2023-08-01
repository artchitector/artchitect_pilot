package main

import (
	"context"
	"github.com/artchitector/artchitect/bot"
	"github.com/artchitector/artchitect/memory"
	"github.com/artchitector/artchitect/model/repository"
	artchitectService "github.com/artchitector/artchitect/soul/core/artchitect"
	artistService "github.com/artchitector/artchitect/soul/core/artist"
	engine2 "github.com/artchitector/artchitect/soul/core/artist/engine"
	"github.com/artchitector/artchitect/soul/core/combinator"
	creator2 "github.com/artchitector/artchitect/soul/core/creator"
	"github.com/artchitector/artchitect/soul/core/entropy"
	"github.com/artchitector/artchitect/soul/core/gifter"
	"github.com/artchitector/artchitect/soul/core/heart"
	"github.com/artchitector/artchitect/soul/core/lottery"
	merciful2 "github.com/artchitector/artchitect/soul/core/merciful"
	"github.com/artchitector/artchitect/soul/core/saver"
	spellerService "github.com/artchitector/artchitect/soul/core/speller"
	"github.com/artchitector/artchitect/soul/core/unifier"
	"github.com/artchitector/artchitect/soul/core/watermark"
	notifier2 "github.com/artchitector/artchitect/soul/notifier"
	"github.com/artchitector/artchitect/soul/resources"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	res := resources.InitResources()
	log.Info().Msg("[main] service soul started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		cancel()
	}()

	// notifier
	notifier := notifier2.NewNotifier(res.GetRedises())

	// Entropy reader + (lightmaster+gatekeeper)
	gk := entropy.NewGatekeeper(res.GetRedises(), notifier)
	lightMaster := entropy.NewLightmaster(res.GetWebcam(), gk)
	go func() {
		if err := lightMaster.StartEntropyReading(ctx); err != nil {
			log.Fatal().Err(err).Msgf("[CRITICAL MALFUNCTION] lightmaster died")
		}
	}()
	entrp := entropy.NewEntropy(lightMaster)

	// repositoties
	artsRepo := repository.NewCardRepository(res.GetDB(), entrp)
	spellRepo := repository.NewSpellRepository(res.GetDB())
	lotteryRepo := repository.NewLotteryRepository(res.GetDB())
	prayRepo := repository.NewPrayRepository(res.GetDB())
	selectionRepo := repository.NewSelectionRepository(res.GetDB())
	unityRepo := repository.NewUnityRepository(res.GetDB())

	// speller+artist+creator
	speller := spellerService.NewSpeller(spellRepo, entrp, notifier)
	var engine artistService.EngineContract
	if res.GetEnv().UseFakeArtist {
		engine = engine2.NewFakeEngine(res.GetEnv().FakeGenerationTime)
	} else {
		engine = engine2.NewArtistEngine(res.GetEnv().ArtistURL)
	}
	sav := saver.NewSaver(res.GetEnv().MemorySaverURL, res.GetEnv().StorageSaverURL)
	watermarkMaker := watermark.NewWatermark()
	artist := artistService.NewArtist(engine, artsRepo, notifier, watermarkMaker, sav)

	// memory (save images to memory-server)
	mmr := memory.NewMemory(res.GetEnv().MemoryHost, nil)

	// artchitect bots
	var artchitectBot *bot.Bot
	if res.GetEnv().Telegram10BotEnabled {
		artchitectBot = bot.NewBot(
			res.GetEnv().Telegram10BotToken,
			artsRepo,
			mmr,
			res.GetEnv().ChatIDArtchitector,
			res.GetEnv().ChatID10,
			res.GetEnv().ChatIDInfinite,
		)
		go artchitectBot.Start(ctx)
	}

	// combinator (makes unity images), unifier (makes unities)
	cmbntr := combinator.NewCombinator(artsRepo, mmr, sav, watermarkMaker)
	unfr := unifier.NewUnifier(unityRepo, artsRepo, entrp, cmbntr, notifier, artchitectBot)

	// creator makes arts
	creator := creator2.NewCreator(
		artist,
		speller,
		notifier,
		unfr,
		artsRepo,
		res.GetEnv().ArtTotalTime,
		res.GetEnv().PrehotDelay,
	)

	// lottery runner
	runner := lottery.NewRunner(lotteryRepo, selectionRepo, artsRepo, entrp, notifier)

	// merciful
	merciful := merciful2.NewMerciful(prayRepo, creator, notifier)

	heartStateOperator := heart.NewHeartState(notifier, artsRepo, 4) // 4 dreams
	go func() {
		if err := heartStateOperator.Run(ctx, 3); err != nil { // 3 seconds
			log.Error().Err(err).Send()
		}
	}()

	// Artchitect core scheduler
	artchitectConfig := artchitectService.Config{
		CardsCreationEnabled: res.GetEnv().CardCreationEnabled,
		LotteryEnabled:       res.GetEnv().LotteryEnabled,
		MercifulEnabled:      res.GetEnv().MercifulEnabled,
		UnifierEnabled:       res.GetEnv().UnifierEnabled,
	}
	artchitect := artchitectService.NewArtchitect(
		artchitectConfig,
		creator,
		lotteryRepo,
		runner,
		merciful,
		unfr,
		notifier,
	)

	// gifter
	if res.GetEnv().GifterActive && artchitectBot != nil {
		gift := gifter.NewGifter(artsRepo, artchitectBot)
		go func() {
			if err := gift.Run(ctx); err != nil {
				log.Fatal().Err(err).Send()
			}
		}()
	}

	//uw := unity_worker.NewUnityWorker(artsRepo, unityRepo)
	//uw.Work(ctx)

	// main loop to make artworks
	var tick int
mainFor:
	for {
		select {
		case <-ctx.Done():
			break mainFor
		case <-time.Tick(time.Second * 1):
			tick += 1
			err := artchitect.Run(ctx, tick)
			if err != nil {
				log.Error().Err(err).Msgf("[main] failed to run artchitect task")
			}
		}
	}

	log.Info().Msg("[main] soul.Setup finished")
}
