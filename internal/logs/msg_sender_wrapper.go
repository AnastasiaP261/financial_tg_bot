package logs

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	"go.uber.org/zap"
)

type msgSender interface {
	SendMessage(text string, userID int64, userName string) error
	SendImage(img []byte, chatID int64, userName string) error
	SendKeyboard(text string, chatID int64, buttonTexts []string, userName string) error

	ListenUpdates(ctx context.Context, msgModel tg.MsgModel)
}

type msgSenderWrapper struct {
	sender msgSender
}

func NewMsgSenderWrapper(origCl msgSender) *msgSenderWrapper {
	return &msgSenderWrapper{
		sender: origCl,
	}
}

func (m msgSenderWrapper) SendMessage(text string, userID int64, userName string) error {
	err := m.sender.SendMessage(text, userID, userName)
	if err != nil {
		logger.Error(
			"send message error",
			zap.Error(err),
			zap.Int64("userId", userID),
		)
		return errors.Wrap(err, "origSender.SendMessage")
	}

	logger.Info(
		"sent messages",
		zap.Int64("userId", userID),
	)

	return nil
}

func (m msgSenderWrapper) SendImage(img []byte, userID int64, userName string) error {
	err := m.sender.SendImage(img, userID, userName)
	if err != nil {
		logger.Error(
			"send image error",
			zap.Error(err),
			zap.Int64("userId", userID),
		)
		return errors.Wrap(err, "origSender.SendImage")
	}

	logger.Info(
		"sent image",
		zap.Int64("userId", userID),
	)

	return nil
}

func (m msgSenderWrapper) SendKeyboard(text string, userID int64, buttonTexts []string, userName string) error {
	err := m.sender.SendKeyboard(text, userID, buttonTexts, userName)
	if err != nil {
		logger.Error(
			"send keyboard error",
			zap.Error(err),
			zap.Int64("userId", userID),
			zap.Strings("button texts", buttonTexts),
		)
		return errors.Wrap(err, "origSender.SendKeyboard")
	}

	logger.Info(
		"sent image",
		zap.Int64("userId", userID),
	)

	return nil
}

func (m msgSenderWrapper) ListenUpdates(ctx context.Context, msgModel tg.MsgModel) {
	logger.Info(
		"messages listen started",
	)

	m.sender.ListenUpdates(ctx, msgModel)
}
