package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
)

type streamer interface {
	GetStream(ctx context.Context, phase string) (chan []byte, func())
}

type StreamRequest struct {
	Phase string `uri:"phase"`
}

type StreamHandler struct {
	streamer streamer
}

func NewStreamHandler(streamer streamer) *StreamHandler {
	return &StreamHandler{streamer}
}

func (sh *StreamHandler) HandleStream(c *gin.Context) {

	r := StreamRequest{}
	if err := c.ShouldBindUri(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	ch, cancel := sh.streamer.GetStream(c, r.Phase)

	c.Header("Content-Type", "multipart/x-mixed-replace;boundary=frame")
	c.Stream(func(w io.Writer) bool {
		multipartWriter := multipart.NewWriter(w)
		multipartWriter.SetBoundary("frame")

		for imgData := range ch {
			iw, err := multipartWriter.CreatePart(textproto.MIMEHeader{
				"Content-type":   []string{"image/jpeg"},
				"Content-length": []string{strconv.Itoa(len(imgData))},
			})
			if err != nil {
				log.Error().Err(err).Msgf("[stream_handler] failed to create part")
				return false
			}
			_, err = iw.Write(imgData)
			if err != nil {
				log.Error().Err(err).Msgf("[stream_handler] failed write image to writer")
				return false
			}
		}
		return false
	})
	log.Info().Msgf("[stream_handler] cancel channel")
	cancel()
}
