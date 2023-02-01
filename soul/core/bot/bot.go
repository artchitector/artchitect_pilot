package bot

import (
	"bytes"
	"context"
	"fmt"
	"github.com/artchitector/artchitect/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

const (
	CommandSendInfinite = "/send_infinite"
)

type cardRepository interface {
	GetCard(ctx context.Context, cardID uint) (model.Card, error)
	GetImage(ctx context.Context, cardID uint) (model.Image, error)
}

type Bot struct {
	token            string
	bot              *bot.Bot
	cardRepository   cardRepository
	artchitectorChat int64
	chat10Min        string
	chatInfinite     string
}

func NewBot(token string, cardRepository cardRepository, artchitectorChat int64, chat10Min string, chatInfinite string) *Bot {
	return &Bot{token, nil, cardRepository, artchitectorChat, chat10Min, chatInfinite}
}

func (t *Bot) Run(ctx context.Context) {
	opts := []bot.Option{
		bot.WithDefaultHandler(t.handler),
	}
	if b, err := bot.New(t.token, opts...); err != nil {
		log.Error().Err(err).Msgf("[bot] failed to create new bot with token %s", t.token)
	} else {
		// handlers
		b.RegisterHandler(bot.HandlerTypeMessageText, CommandSendInfinite, bot.MatchTypePrefix, t.infiniteHandler)
		// start bot to listen all messages
		log.Info().Msgf("[bot] starting bot")
		t.bot = b
		b.Start(ctx)
		log.Info().Msgf("[bot] bot finished")
	}
}

func (t *Bot) SendCardTo10Min(ctx context.Context, cardID uint) error {
	return t.sendCard(ctx, cardID, t.chat10Min)
}

func (t *Bot) sendCardToInfinite(ctx context.Context, cardID uint) error {
	return t.sendCard(ctx, cardID, t.chatInfinite)
}

func (t *Bot) sendCard(ctx context.Context, cardID uint, chatID string) error {
	if t.bot == nil {
		return errors.Errorf("[bot] not initialized")
	}
	card, err := t.getCard(ctx, cardID)
	if err != nil {
		return errors.Wrapf(err, "[bot] failed to get card by ID=%d", cardID)
	}
	if card.Image.CardID == 0 {
		return errors.Errorf("[bot] not found image in card(id=%d), given to bot", card.ID)
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
	r := bytes.NewReader(card.Image.Data)
	msg, err := t.bot.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID:  chatID,
		Photo:   &models.InputFileUpload{Data: r},
		Caption: text,
	})
	if err != nil {
		return errors.Wrapf(err, "[bot] failed send photo for card id=%d", card.ID)
	}

	log.Info().Msgf("[bot] sent card to chat %s. CardID=%d. MessageID=%d", chatID, card.ID, msg.ID)
	return nil
}

func (t *Bot) getCard(ctx context.Context, cardID uint) (model.Card, error) {
	card, err := t.cardRepository.GetCard(ctx, cardID)
	if err != nil {
		return model.Card{}, errors.Wrapf(err, "[bot] failed to GetCard %d", cardID)
	}
	image, err := t.cardRepository.GetImage(ctx, card.ID)
	if err != nil {
		return model.Card{}, errors.Wrapf(err, "[bot] failed to get image for card %d", card.ID)
	} else {
		card.Image = image
	}
	return card, nil
}

func (t *Bot) handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.ChannelPost != nil {
		log.Info().Msgf("[bot_infinite] got channel post with default handler: %+v", update.ChannelPost)
	} else {
		log.Info().Msgf("[bot_infinite] got message with default handler: %+v", update.Message)
	}
}

func (t *Bot) infiniteHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if err := t.checkArtchitector(update.Message); err != nil {
		log.Error().Err(err).Msgf("[bot_infinite] security check failed to handle infinite")
		return
	}
	args := t.parseArguments(update.Message.Text)
	if len(args) == 0 {
		t.replyError(ctx, update.Message, errors.Errorf("[bot_infinite] need cardID as argument"))
		return
	}
	cardID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		t.replyError(ctx, update.Message, errors.Errorf("[bot_infinite] cardID must be uint"))
		return
	}
	log.Info().Msgf("[bot_infinite] got infinite command with cardID", cardID)
	err = t.sendCardToInfinite(ctx, uint(cardID))
	if err != nil {
		t.replyError(ctx, update.Message, err)
	}
}

func (t *Bot) checkArtchitector(msg *models.Message) error {
	if msg.Chat.ID != t.artchitectorChat {
		return errors.Errorf("[bot] not artchitector chat. ChatID: %d", msg.Chat.ID)
	}
	return nil
}

func (t *Bot) parseArguments(text string) []string {
	args := strings.Split(text, " ")
	return args[1:]
}

func (t *Bot) replyError(ctx context.Context, msg *models.Message, errMessage error) {
	_, err := t.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text:   errMessage.Error(),
	})
	if err != nil {
		log.Error().Err(err).Msg("[bot_infinite] failed to reply error")
	}
}
