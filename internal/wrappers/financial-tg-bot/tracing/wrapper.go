package logs

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/financial-tg-bot"
	"go.uber.org/zap"
)

func InitTracing(logger *zap.Logger, serviceName string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}
	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}
}

type Wrapper struct {
	sender financial_tg_bot.MsgSender
}

func NewWrapper(origCl financial_tg_bot.MsgSender) *Wrapper {
	return &Wrapper{
		sender: origCl,
	}
}

func (m *Wrapper) SendMessage(text string, userID int64) error {
	return m.sender.SendMessage(text, userID)
}

func (m *Wrapper) SendImage(img []byte, userID int64) error {
	return m.sender.SendImage(img, userID)
}

func (m *Wrapper) SendKeyboard(text string, userID int64, buttonTexts []string) error {
	return m.sender.SendKeyboard(text, userID, buttonTexts)
}

func (m *Wrapper) IncomingCallback(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "incoming callback")
	defer span.Finish()

	return m.sender.IncomingCallback(ctx, model, msg)
}

func (m *Wrapper) IncomingMessage(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "incoming message")
	defer span.Finish()

	return m.sender.IncomingMessage(ctx, model, msg)
}
