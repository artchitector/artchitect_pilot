package handler

import (
	"encoding/json"
	"github.com/artchitector/artchitect/gate/localmodel"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebsocketHandler struct {
	listener listener
	eventsCh chan localmodel.Event
}

func NewWebsocketHandler(l listener) *WebsocketHandler {
	return &WebsocketHandler{listener: l}
}

func (wh *WebsocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Error().Err(err).Msgf("[websocket] failed to Upgrade")
		return
	}

	ch, done := wh.listener.EventChannel()
	defer close(ch)

	for event := range ch {
		j, err := json.Marshal(event)
		if err != nil {
			log.Error().Err(err).Msgf("[websocket] failed to marshal event")
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, j); err != nil {
			log.Warn().Err(err).Msgf("[websocket] failed to write event")
			done <- struct{}{}
			return
		}
	}

	log.Warn().Msgf("[websocket] stopping handler")
}
