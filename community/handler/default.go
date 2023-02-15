package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/rs/zerolog/log"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Info().Msgf("[default_handler] start")
	if update.ChannelPost != nil {
		log.Warn().Msgf("[default_handler] got message from channel: %d %s", update.ChannelPost.ID, update.Message.Text)
		return
	}

	var message *models.Message
	if update.Message != nil {
		message = update.Message
	} else if update.EditedMessage != nil {
		message = update.EditedMessage
	}
	if message == nil {
		log.Error().Msgf("[default_handler] nil message")
		return
	}
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           message.Chat.ID,
		Text:             fmt.Sprintf("*unknown command*: %s", message.Text),
		ParseMode:        models.ParseModeMarkdown,
		ReplyToMessageID: message.ID,
	})
	if err != nil {
		log.Error().Err(err).Msgf("[default_handler] failed to send reply message")
	}
	log.Info().Msgf("[default_handler] sent message back to chat %d", message.Chat.ID)
}
