package logs

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"go.uber.org/zap"
)

type msgSender interface {
	SendMessage(text string, userID int64) error
	SendImage(img []byte, userID int64) error
	SendKeyboard(text string, userID int64, buttonTexts []string) error

	IncomingCallback(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error
	IncomingMessage(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error
}

type MsgSenderWrapper struct {
	sender msgSender
}

func NewMsgSenderWrapper(origCl msgSender) *MsgSenderWrapper {
	return &MsgSenderWrapper{
		sender: origCl,
	}
}

func (m *MsgSenderWrapper) SendMessage(text string, userID int64) error {
	err := m.sender.SendMessage(text, userID)
	if err != nil {
		logs.Error(
			"send message error",
			zap.Error(err),
			zap.Int64("userId", userID),
		)
		return errors.Wrap(err, "origSender.SendMessage")
	}

	logs.Info(
		"sent messages",
		zap.Int64("userId", userID),
	)

	return nil
}

func (m *MsgSenderWrapper) SendImage(img []byte, userID int64) error {
	err := m.sender.SendImage(img, userID)
	if err != nil {
		logs.Error(
			"send image error",
			zap.Error(err),
			zap.Int64("userId", userID),
		)
		return errors.Wrap(err, "origSender.SendImage")
	}

	logs.Info(
		"sent image",
		zap.Int64("userId", userID),
	)

	return nil
}

func (m *MsgSenderWrapper) SendKeyboard(text string, userID int64, buttonTexts []string) error {
	err := m.sender.SendKeyboard(text, userID, buttonTexts)
	if err != nil {
		logs.Error(
			"send keyboard error",
			zap.Error(err),
			zap.Int64("userId", userID),
			zap.Strings("button texts", buttonTexts),
		)
		return errors.Wrap(err, "origSender.SendKeyboard")
	}

	logs.Info(
		"sent image",
		zap.Int64("userId", userID),
	)

	return nil
}

func (m *MsgSenderWrapper) IncomingCallback(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	err := m.sender.IncomingCallback(ctx, model, msg)
	if err != nil {
		logs.Error(
			"income message error",
			zap.Error(err),
			zap.Int64("userId", msg.CallbackQuery.From.ID),
		)
		return errors.Wrap(err, "sender.IncomingCallback")
	}

	logs.Info(
		"incoming messages",
		zap.Int64("userId", msg.CallbackQuery.From.ID),
	)

	return nil
}

func (m *MsgSenderWrapper) IncomingMessage(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	if err := m.sender.IncomingMessage(ctx, model, msg); err != nil {
		logs.Error(
			"income message error",
			zap.Error(err),
			zap.Int64("userId", msg.Message.From.ID),
		)
		return errors.Wrap(err, "sender.IncomingMessage")
	}

	logs.Info(
		"incoming messages",
		zap.Int64("userId", msg.Message.From.ID),
	)

	return nil
}
