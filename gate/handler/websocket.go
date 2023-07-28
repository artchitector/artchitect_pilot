package handler

import (
	"encoding/json"
	"fmt"
	"github.com/artchitector/artchitect/gate/localmodel"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
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

	subscribedChannels := []string{}
	go func() {
		// read messages. channel subscription goes from client
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error().Err(err).Msgf("[websocket] got error from conn.ReadMessage")
				} else {
					log.Warn().Msgf("[websocket] client closed correctly")
				}
				done <- struct{}{}
				break
			}

			subscribedChannels, err = wh.checkAndSubscribe(subscribedChannels, string(message))
			if err != nil {
				log.Error().Err(err).Send()
				break
			}
		}
	}()

	for event := range ch {
		// check, if we have subscription in this ws-channel
		if !wh.isSubscribed(subscribedChannels, event.Name) {
			continue
		}

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

func (wh *WebsocketHandler) isSubscribed(channels []string, name string) bool {
	for _, ch := range channels {
		if name == ch {
			return true
		}
	}
	return false
}

func (wh *WebsocketHandler) checkAndSubscribe(channels []string, s string) ([]string, error) {
	if wh.isSubscribed(channels, s) {
		return channels, nil
	}
	segments := strings.Split(s, ".")
	if len(segments) != 2 || segments[0] != "subscribe" {
		return channels, fmt.Errorf("[websocket] failed to parse message: %s", s)
	}
	channels = append(channels, segments[1])
	log.Info().Msgf("[websocket] subscribe to channel %s", segments[1])
	return channels, nil
}
