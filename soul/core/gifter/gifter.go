package gifter

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type cardRepository interface {
	GetTotalCards(ctx context.Context) (uint, error)
	GetCardWithOffset(offset uint) (model.Card, error)
}

type origin interface {
	Select(ctx context.Context, totalVariants uint) (uint, error)
}

const (
	MaxAttempts = 10 //each photo 10 times because of context timeout
)

type Gifter struct {
	origin           origin
	cardRepository   cardRepository
	telegramBotToken string
	tenMinChatID     string
}

func NewGifter(
	cardRepository cardRepository,
	origin origin,
	telegramBotToken string,
	tenMinChatID string,
) *Gifter {
	return &Gifter{origin, cardRepository, telegramBotToken, tenMinChatID}
}

func (g *Gifter) Run(ctx context.Context) error {
	for {
		currentAttempts := 0
		for {
			currentAttempts += 1
			if currentAttempts > MaxAttempts {
				log.Info().Msgf("[gifter] max attempts (%d) exceeded", MaxAttempts)
				break
			}
			err := g.sendCard(ctx)
			if err != nil {
				log.Error().Err(err).Msgf("[gifter] failed to send card")
			} else {
				break
			}
		}
		time.Sleep(time.Minute * 10)
	}
}

func (g *Gifter) sendCard(ctx context.Context) error {
	log.Info().Msgf("[gifter] !!SEND GIFT!!")
	card, err := g.getCard(ctx)
	if err != nil {
		return errors.Wrap(err, "[gifter] failed to getCard")
	}
	b, err := g.getBot()
	if err != nil {
		return errors.Wrap(err, "[gifter] failed to getBot")
	}

	text := fmt.Sprintf(
		"Card #%d. (https://artchitect.space/card/%d)\n\n"+
			"Created: %s\n"+
			"Seed: %d\n"+
			"Tags: %s",
		card.ID,
		card.ID,
		card.CreatedAt.Format("2006 Jan 2 15:04"),
		card.Spell.Seed,
		card.Spell.Tags,
	)

	r := bytes.NewReader(card.Image)
	msg, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:    g.tenMinChatID,
		Photo:     &models.InputFileUpload{Data: r},
		Caption:   text,
		ParseMode: "",
	})
	if err != nil {
		return errors.Wrap(err, "[gifter] failed send photo")
	}

	log.Info().Msgf("[gifter] sent gift. CardID=%d. MessageID=%d", card.ID, msg.ID)
	return nil
}

func (g *Gifter) getCard(ctx context.Context) (model.Card, error) {
	totalCards, err := g.cardRepository.GetTotalCards(ctx)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[gifter] failed get total cards")
	}
	selection, err := g.origin.Select(ctx, totalCards)
	if err != nil {
		return model.Card{}, errors.Wrap(err, "[gifter] failed to select from origin")
	}
	card, err := g.cardRepository.GetCardWithOffset(selection)
	if err != nil {
		return model.Card{}, errors.Wrapf(err, "[gifter] failed to GetCardWithOffset %d", selection-1)
	}
	return card, nil
}

func (g *Gifter) getBot() (*bot.Bot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}
	return bot.New(g.telegramBotToken, opts...)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Info().Msgf("[gifter] incoming message. ChatID: %d, Text: %s", update.Message.Chat.ID, update.Message.Text)
}
