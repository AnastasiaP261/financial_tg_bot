package logs

import (
	"context"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"go.uber.org/zap"
)

type Wrapper struct {
	sender wrappers.MsgSender
}

func NewWrapper(origCl wrappers.MsgSender) *Wrapper {
	return &Wrapper{
		sender: origCl,
	}
}

func (m *Wrapper) SendMessage(text string, userID int64) error {
	err := m.sender.SendMessage(text, userID)
	if err != nil {
		logs.Error(
			"send message error",
			zap.Error(err),
			zap.Int64("userId", userID),
		)
		return err
	}

	logs.Info(
		"sent messages",
		zap.Int64("userId", userID),
	)

	return nil
}

func (m *Wrapper) SendImage(img []byte, userID int64) error {
	err := m.sender.SendImage(img, userID)
	if err != nil {
		logs.Error(
			"send image error",
			zap.Error(err),
			zap.Int64("userId", userID),
		)
		return err
	}

	logs.Info(
		"sent image",
		zap.Int64("userId", userID),
	)

	return nil
}

func (m *Wrapper) SendKeyboard(text string, userID int64, buttonTexts []string) error {
	err := m.sender.SendKeyboard(text, userID, buttonTexts)
	if err != nil {
		logs.Error(
			"send keyboard error",
			zap.Error(err),
			zap.Int64("userId", userID),
			zap.Strings("button texts", buttonTexts),
		)
		return err
	}

	logs.Info(
		"sent image",
		zap.Int64("userId", userID),
	)

	return nil
}

func (m *Wrapper) IncomingCallback(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	err := m.sender.IncomingCallback(ctx, model, msg)
	if err != nil {
		logs.Error(
			"income message error",
			zap.Error(err),
			zap.Int64("userId", msg.CallbackQuery.From.ID),
		)
		return err
	}

	logs.Info(
		"incoming messages",
		zap.Int64("userId", msg.CallbackQuery.From.ID),
	)

	return nil
}

func (m *Wrapper) IncomingMessage(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	if err := m.sender.IncomingMessage(ctx, model, msg); err != nil {
		logs.Error(
			"income message error",
			zap.Error(err),
			zap.Int64("userId", msg.Message.From.ID),
		)
		return err
	}

	logs.Info(
		"incoming messages",
		zap.Int64("userId", msg.Message.From.ID),
	)

	return nil
}
