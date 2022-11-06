package logs

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/messages"
	"go.uber.org/zap"
)

type msgReceiver interface {
	IncomingCallback(ctx context.Context, msg messages.Callback) error
	IncomingMessage(ctx context.Context, msg messages.Message) error
}

type msgReceiverWrapper struct {
	receiver msgReceiver
}

func NewMsgReceiverWrapper(origCl msgReceiver) *msgReceiverWrapper {
	return &msgReceiverWrapper{
		receiver: origCl,
	}
}

func (w *msgReceiverWrapper) IncomingMessage(ctx context.Context, msg messages.Message) error {
	err := w.receiver.IncomingMessage(ctx, msg)
	if err != nil {
		logger.Error(
			"income message error",
			zap.Error(err),
			zap.Int64("userId", msg.UserID),
		)
		return errors.Wrap(err, "receiver.IncomingMessage")
	}

	logger.Info(
		"incoming messages",
		zap.Int64("userId", msg.UserID),
	)

	return nil
}

func (w *msgReceiverWrapper) IncomingCallback(ctx context.Context, msg messages.Callback) error {
	err := w.receiver.IncomingCallback(ctx, msg)
	if err != nil {
		logger.Error(
			"income message error",
			zap.Error(err),
			zap.Int64("userId", msg.UserID),
		)
		return errors.Wrap(err, "receiver.IncomingCallback")
	}

	logger.Info(
		"incoming messages",
		zap.Int64("userId", msg.UserID),
	)

	return nil
}
