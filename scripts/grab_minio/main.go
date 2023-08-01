package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"math"
	"net/http"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"})

	for {
		if lastCard, err := grabNextCard(); err != nil {
			if lastCard >= 278428 { // 278428 {
				return
			}
			log.Error().Err(err).Msgf("[grab_minio] failed to get card %d", lastCard)
		} else {
			if lastCard >= 278428 { // 278428 {
				return
			}
			log.Info().Msgf("[grab_minio] finished card %d", lastCard)
		}
	}
}

func getLastCard() int {
	b, err := os.ReadFile("./last_card")
	if err != nil {
		log.Fatal().Err(err)
	}
	var lastCard int
	if err := json.Unmarshal(b, &lastCard); err != nil {
		log.Fatal().Err(err)
	}
	return lastCard
}

func saveLastCard(lastCard int) error {
	return os.WriteFile("./last_card", []byte(fmt.Sprintf("%d", lastCard)), 0777)
}

func grabNextCard() (int, error) {
	lastCard := getLastCard()
	lastCard += 1
	if err := saveLastCard(lastCard); err != nil {
		log.Fatal().Err(err)
	}

	url := fmt.Sprintf("http://storage.artchitect.space/cards/card-%d.jpg", lastCard)
	resp, err := http.Get(url)
	if err != nil {
		return lastCard, err
	} else if resp.StatusCode != http.StatusOK {
		return lastCard, errors.Errorf("bad status %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return lastCard, err
	}

	kNumber := math.Floor(float64(lastCard) / 10000.0)
	path := fmt.Sprintf("/home/artchitector/storage_dump/storage_cards/%d", int(kNumber))

	if err := os.MkdirAll(path, 0777); err != nil {
		return 0, err
	}

	fileName := fmt.Sprintf("%s/%d.jpg", path, lastCard)
	if err := os.WriteFile(fileName, data, 0777); err != nil {
		return lastCard, err
	}

	log.Info().Msgf("[grab_minio] saved card %d into %s", lastCard, fileName)

	return lastCard, nil
}
