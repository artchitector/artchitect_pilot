package bot

import (
	"bytes"
	"context"
	"github.com/artchitector/artchitect/model"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

const (
	CommandSendInfinite = "/send_infinite"
	CommandGive         = "/give" // give one selected card from all
)

type cardRepository interface {
	GetCard(ctx context.Context, cardID uint) (model.Card, error)
	GetOriginSelectedCard(ctx context.Context) (model.Card, error)
	GetOriginSelectedCardByPeriod(ctx context.Context, start time.Time, end time.Time) (model.Card, error)
}

type memory interface {
	DownloadImage(ctx context.Context, cardID uint, size string) ([]byte, error)
}

type Bot struct {
	token            string
	bot              *bot.Bot
	cardRepository   cardRepository
	memory           memory
	artchitectorChat int64
	chat10Min        string
	chatInfinite     string
}

func NewBot(token string, cardRepository cardRepository, memory memory, artchitectorChat int64, chat10Min string, chatInfinite string) *Bot {
	return &Bot{token, nil, cardRepository, memory, artchitectorChat, chat10Min, chatInfinite}
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
		b.RegisterHandler(bot.HandlerTypeMessageText, CommandGive, bot.MatchTypePrefix, t.giveHandler)
		// start bot to listen all messages
		log.Info().Msgf("[bot] starting bot")
		t.bot = b
		b.Start(ctx)
		log.Info().Msgf("[bot] bot finished")
	}
}

func (t *Bot) SendCardTo10Min(ctx context.Context, cardID uint) error {
	card, img, err := t.getCard(ctx, cardID)
	if err != nil {
		return errors.Wrapf(err, "[bot] failed to get card by ID=%d", cardID)
	}
	text := getTextWithoutCaption(card)
	return t.sendCard(ctx, card, img, text, t.chat10Min)
}

func (t *Bot) sendCardToInfinite(ctx context.Context, cardID uint, caption string) error {
	card, img, err := t.getCard(ctx, cardID)
	if err != nil {
		return errors.Wrapf(err, "[bot] failed to get card by ID=%d", cardID)
	}
	var text string
	if caption == "" {
		text = getTextWithoutCaption(card)
	} else {
		text = getTextWithCaption(card, caption)
	}
	return t.sendCard(ctx, card, img, text, t.chatInfinite)
}

func (t *Bot) sendCardBack(ctx context.Context, cardID uint, chatID string) error {
	card, img, err := t.getCard(ctx, cardID)
	if err != nil {
		return errors.Wrapf(err, "[bot] failed to get card by ID=%d", cardID)
	}
	return t.sendCard(ctx, card, img, getTextWithoutCaption(card), chatID)
}

func (t *Bot) sendCard(ctx context.Context, card model.Card, img []byte, text string, chatID string) error {
	if t.bot == nil {
		return errors.Errorf("[bot] not initialized")
	}

	r := bytes.NewReader(img)
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

func (t *Bot) getCard(ctx context.Context, cardID uint) (model.Card, []byte, error) {
	card, err := t.cardRepository.GetCard(ctx, cardID)
	if err != nil {
		return model.Card{}, nil, errors.Wrapf(err, "[bot] failed to GetCard %d", cardID)
	}

	image, err := t.memory.DownloadImage(ctx, cardID, model.SizeF)
	if err != nil {
		return model.Card{}, nil, errors.Wrapf(err, "[bot] failed to get image for card %d", card.ID)
	}
	return card, image, nil
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

	var caption string
	if len(args) > 1 {
		caption = strings.Join(args[1:], " ")
	}

	log.Info().Msgf("[bot_infinite] got infinite command with cardID", cardID)
	err = t.sendCardToInfinite(ctx, uint(cardID), caption)
	if err != nil {
		t.replyError(ctx, update.Message, err)
	}
}

func (t *Bot) giveHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		t.replyError(ctx, update.Message, errors.Errorf("only in private chats available"))
		return
	}

	var card model.Card
	var err error
	args := t.parseArguments(update.Message.Text)
	if len(args) == 0 {
		card, err = t.cardRepository.GetOriginSelectedCard(ctx)
		if err != nil {
			log.Error().Err(err).Msgf("[bot_give] failed to get card")
			t.replyError(ctx, update.Message, errors.Errorf("failed to get card. try once more"))
			return
		}
	} else if len(args) == 1 {
		dur, err := time.ParseDuration(args[0])
		if err != nil {
			t.replyError(ctx, update.Message, errors.Errorf("failed to parse duration in string %s", args[0]))
			return
		}
		now := time.Now()
		start := now.Add(-1 * dur)
		card, err = t.cardRepository.GetOriginSelectedCardByPeriod(ctx, start, now)
		if err != nil {
			log.Error().Err(err).Msgf("[bot_give] failed to get GetOriginSelectedCardByPeriod")
			t.replyError(ctx, update.Message, errors.Errorf("failed to get card with period. try once more"))
			return
		}
	} else {
		t.replyError(ctx, update.Message, errors.Errorf("[bot_infinite] too many arguments"))
		return
	}

	if err := t.sendCardBack(ctx, card.ID, strconv.FormatInt(update.Message.Chat.ID, 10)); err != nil {
		log.Error().Err(err).Msgf("[bot_give] failed to send card %d", card.ID)
		t.replyError(ctx, update.Message, errors.Errorf("failed to send message. try once more"))
		return
	}
	log.Info().Msgf("[bot_give] sent given card %d to chat %d", card.ID, update.Message.Chat.ID)
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
