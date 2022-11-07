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

type MsgModel interface {
	IncomingCallback(ctx context.Context, msg messages.Callback) error
	IncomingMessage(ctx context.Context, msg messages.Message) error
}

type MsgHandler interface {
	SendMessage(text string, userID int64) error
	SendImage(img []byte, userID int64) error
	SendKeyboard(text string, userID int64, buttonTexts []string) error

	IncomingCallback(ctx context.Context, model MsgModel, msg tgbotapi.Update) error
	IncomingMessage(ctx context.Context, model MsgModel, msg tgbotapi.Update) error
}

type Client struct {
	client     *tgbotapi.BotAPI
	msgHandler MsgHandler
}

func New(tokenGetter TokenGetter, msgHandler MsgHandler) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client:     client,
		msgHandler: msgHandler,
	}, nil
}

func (c *Client) ListenUpdates(ctx context.Context, msgModel MsgModel) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			_ = c.msgHandler.IncomingCallback(ctx, msgModel, update)

		case update.Message != nil:
			_ = c.msgHandler.IncomingMessage(ctx, msgModel, update)
		}
	}
}
