package main

import (
	"context"
	"github.com/artchitector/artchitect/saver/handler"
	"github.com/artchitector/artchitect/saver/resources"
	"github.com/artchitector/artchitect/saver/saver"
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
    - all images are in /var/artchitect/arts folder (set in env)
    - 10k-arts is in separate folder folder=(id % 10000)
    - arts names in these folders:
    art-56910-f.jpg
    art-56910-m.jpg
    art-56910-s.jpg
    art-56910-xs.jpg
    these files statically served by nginx, and gate services can take img and proxy it
*/
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})

	res := resources.InitResources()
	log.Info().Msg("service gate started")
	svr := saver.NewSaver(res.GetEnv().ArtsPath, res.GetEnv().UnityPath, res.GetEnv().FullSizePath)
	uploadHandler := handler.NewUploadHandler(svr, res.GetEnv().IsFullsizeStorage)

	go func() {
		r := gin.Default()
		r.MaxMultipartMemory = 8 << 20 // 8 MiB
		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
		}))
		if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
			log.Fatal().Err(err).Send()
		}
		r.POST("/upload_art", uploadHandler.Handle)
		r.POST("/upload_unity", uploadHandler.HandleUnity)
		r.POST("/upload_fullsize", uploadHandler.HandleFullsize)
		if err := r.Run("0.0.0.0:" + res.GetEnv().HttpPort); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("saver.Setup finished")
}
