package bot

import (
	"context"
	"github.com/artchitector/artchitect/community/handler"
	"github.com/go-telegram/bot"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	CommandLogin = "/login"
)

type Bot struct {
	token string
	bot   *bot.Bot
}

func NewBot(token string) *Bot {
	return &Bot{token, nil}
}

func (t *Bot) Run(ctx context.Context) error {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler.DefaultHandler),
	}
	if b, err := bot.New(t.token, opts...); err != nil {
		return errors.Wrap(err, "[bot] failed to run bot")
	} else {
		log.Info().Msgf("[bot] starting bot")
		t.bot = b
		b.Start(ctx)
		log.Info().Msgf("[bot] bot finished")
	}
	return nil
}
