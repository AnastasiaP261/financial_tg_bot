package tg

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/messages"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client *tgbotapi.BotAPI
}

type MsgModel interface {
	IncomingCallback(ctx context.Context, msg messages.Callback) error
	IncomingMessage(ctx context.Context, msg messages.Message) error
}

func New(tokenGetter TokenGetter) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SendMessage(text string, userID int64, userName string) error {
	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}

	return nil
}

func (c *Client) SendImage(img []byte, userId int64, userName string) error {
	b := tgbotapi.FileBytes{Bytes: img}

	_, err := c.client.Send(tgbotapi.NewPhoto(userId, b))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}

	return nil
}

func (c *Client) SendKeyboard(text string, userId int64, buttonTexts []string, userName string) error {
	msg := tgbotapi.NewMessage(userId, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, txt := range buttonTexts {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(txt, txt)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard

	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}

	return nil
}

func (c *Client) ListenUpdates(ctx context.Context, msgModel MsgModel) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			_ = msgModel.IncomingCallback(ctx, messages.Callback{
				UserID:   update.CallbackQuery.From.ID,
				UserName: update.CallbackQuery.From.UserName,
				Data:     update.CallbackQuery.Data,
			})

		case update.Message != nil:
			_ = msgModel.IncomingMessage(ctx, messages.Message{
				Text:     update.Message.Text,
				UserID:   update.Message.From.ID,
				UserName: update.Message.From.UserName,
			})
		}
	}
}
