package main

import (
	"context"
	"github.com/artchitector/artchitect/saver/handler"
	"github.com/artchitector/artchitect/saver/resources"
	"github.com/artchitector/artchitect/saver/saver"
	"github.com/artchitector/artchitect/saver/worker"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

/*
 1. Saver service saves image in different sized in file system
    file structure:
    - all images are in /root/cards folder
    - every 10k cards is in separate folder folder=(id % 10000)
    - card names in these folders:
    card-56910-f.jpg
    card-56910-m.jpg
    card-56910-s.jpg
    card-56910-xs.jpg
    these files statically served by nginx, and gate services can take img and proxy it
 2. NO DIRECT CONNECTIONS FROM CLIENTS.
    Sometimes saver will be under VPN with only internal access
    Gate take jpg and proxy it to the client
 3. No image storage inside database. Database is lightweight - not binaries.
    ID is universal to get any file
 4. Initial scripts to download images from database and put it as files (one time script)
*/
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})

	res := resources.InitResources()
	log.Info().Msg("service gate started")
	svr := saver.NewSaver(res.GetEnv().CardsPath)
	uploadHandler := handler.NewUploadHandler(svr)
	wrk := worker.NewWorker(res.GetDB(), svr)
	go func() {
		wrk.Work(ctx)
	}()

	go func() {
		r := gin.Default()
		r.MaxMultipartMemory = 8 << 20 // 8 MiB
		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
		}))
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal().Err(err).Send()
		}
		r.POST("/upload", uploadHandler.Handle)
		if err := r.Run("0.0.0.0:" + res.GetEnv().HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("saver.Run finished")
}
