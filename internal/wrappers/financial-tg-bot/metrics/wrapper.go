package metrics

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/utils/metrics"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/financial-tg-bot"
)

type Wrapper struct {
	sender financial_tg_bot.MsgSender
}

func NewWrapper(origCl financial_tg_bot.MsgSender) *Wrapper {
	return &Wrapper{
		sender: origCl,
	}
}

func (m *Wrapper) SendMessage(text string, userID int64) error {
	err := m.sender.SendMessage(text, userID)
	if err != nil {
		metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeOutgoing, metrics.StatusErr).Inc()
		return err
	}

	metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeOutgoing, metrics.StatusOk).Inc()

	return nil
}

func (m *Wrapper) SendImage(img []byte, userID int64) error {
	err := m.sender.SendImage(img, userID)
	if err != nil {
		metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeOutgoing, metrics.StatusErr).Inc()
		return err
	}

	metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeOutgoing, metrics.StatusOk).Inc()

	return nil
}

func (m *Wrapper) SendKeyboard(text string, userID int64, buttonTexts []string) error {
	err := m.sender.SendKeyboard(text, userID, buttonTexts)
	if err != nil {
		metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeOutgoing, metrics.StatusErr).Inc()
		return err
	}

	metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeOutgoing, metrics.StatusOk).Inc()

	return nil
}

func (m *Wrapper) IncomingCallback(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	startTime := time.Now()
	err := m.sender.IncomingCallback(ctx, model, msg)
	if err != nil {
		metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeIncoming, metrics.StatusErr).Inc()
		return err
	}
	duration := time.Since(startTime)

	metrics.SummaryResponseTime.Observe(duration.Seconds())
	metrics.HistogramResponseTime.Observe(duration.Seconds())

	metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeIncoming, metrics.StatusOk).Inc()

	return nil
}

func (m *Wrapper) IncomingMessage(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	startTime := time.Now()
	if err := m.sender.IncomingMessage(ctx, model, msg); err != nil {
		metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeIncoming, metrics.StatusErr).Inc()
		return err
	}
	duration := time.Since(startTime)

	metrics.SummaryResponseTime.Observe(duration.Seconds())
	metrics.HistogramResponseTime.Observe(duration.Seconds())

	metrics.InFlightTypeMsg.WithLabelValues(metrics.TypeIncoming, metrics.StatusOk).Inc()

	return nil
}
