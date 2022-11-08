package wrappers

import (
	"context"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
)

type MsgSender interface {
	SendMessage(text string, userID int64) error
	SendImage(img []byte, userID int64) error
	SendKeyboard(text string, userID int64, buttonTexts []string) error

	IncomingCallback(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error
	IncomingMessage(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error
}
