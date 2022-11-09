package metrics

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers"
)

const (
	TypeIncoming = "incoming"
	TypeOutgoing = "outgoing"

	StatusOk  = "ok"
	StatusErr = "error"

	ReportSourceBD    = "bd"
	ReportSourceCache = "cache"
)

var (
	// InFlightTypeMsg количество входящих и исходящих по статусам
	InFlightTypeMsg = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "tg_bot",
			Subsystem: "msg",
			Name:      "in_flight_msg_for_type",
		},
		[]string{"type", "status"},
	)

	// InFlightCommandMsg количество входящих по командам и статусам
	InFlightCommandMsg = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "tg_bot",
			Subsystem: "msg",
			Name:      "in_flight_msg_for_command",
		},
		[]string{"command", "status"},
	)

	// SummaryResponseTime саммари времени ответа бота
	SummaryResponseTime = promauto.NewSummary(prometheus.SummaryOpts{
		Namespace: "tg_bot",
		Subsystem: "msg",
		Name:      "summary_response_time_seconds",
		Objectives: map[float64]float64{
			0.5:  0.1,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	// HistogramResponseTime гистограмма времени ответа бота
	HistogramResponseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "tg_bot",
		Subsystem: "msg",
		Name:      "histogram_response_time_seconds",
		Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2},
		// Buckets: prometheus.ExponentialBucketsRange(0.0001, 2, 16),
	})

	// SummaryFixerResponseTime саммари времени ответа fixer
	SummaryFixerResponseTime = promauto.NewSummary(prometheus.SummaryOpts{
		Namespace: "tg_bot",
		Subsystem: "fixer",
		Name:      "summary_fixer_response_time_seconds",
		Objectives: map[float64]float64{
			0.5:  0.1,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	// HistogramFixerResponseTime гистограмма времени ответа fixer
	HistogramFixerResponseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "tg_bot",
		Subsystem: "fixer",
		Name:      "histogram_fixer_response_time_seconds",
		Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		// Buckets: prometheus.ExponentialBucketsRange(0.0001, 2, 16),
	})

	// InFlightCache количество успешных и с ошибкой ответов кеша
	InFlightCache = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "tg_bot",
			Subsystem: "msg",
			Name:      "in_flight_cache",
		},
		[]string{"status"},
	)

	// InFlightReports количество успешно сформированных отчетов (по источнику)
	InFlightReports = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "tg_bot",
			Subsystem: "msg",
			Name:      "in_flight_reports",
		},
		[]string{"source"},
	)
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
		InFlightTypeMsg.WithLabelValues(TypeOutgoing, StatusErr).Inc()
		return err
	}

	InFlightTypeMsg.WithLabelValues(TypeOutgoing, StatusOk).Inc()

	return nil
}

func (m *Wrapper) SendImage(img []byte, userID int64) error {
	err := m.sender.SendImage(img, userID)
	if err != nil {
		InFlightTypeMsg.WithLabelValues(TypeOutgoing, StatusErr).Inc()
		return err
	}

	InFlightTypeMsg.WithLabelValues(TypeOutgoing, StatusOk).Inc()

	return nil
}

func (m *Wrapper) SendKeyboard(text string, userID int64, buttonTexts []string) error {
	err := m.sender.SendKeyboard(text, userID, buttonTexts)
	if err != nil {
		InFlightTypeMsg.WithLabelValues(TypeOutgoing, StatusErr).Inc()
		return err
	}

	InFlightTypeMsg.WithLabelValues(TypeOutgoing, StatusOk).Inc()

	return nil
}

func (m *Wrapper) IncomingCallback(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	startTime := time.Now()
	err := m.sender.IncomingCallback(ctx, model, msg)
	if err != nil {
		InFlightTypeMsg.WithLabelValues(TypeIncoming, StatusErr).Inc()
		return err
	}
	duration := time.Since(startTime)

	SummaryResponseTime.Observe(duration.Seconds())
	HistogramResponseTime.Observe(duration.Seconds())

	InFlightTypeMsg.WithLabelValues(TypeIncoming, StatusOk).Inc()

	return nil
}

func (m *Wrapper) IncomingMessage(ctx context.Context, model tg.MsgModel, msg tgbotapi.Update) error {
	startTime := time.Now()
	if err := m.sender.IncomingMessage(ctx, model, msg); err != nil {
		InFlightTypeMsg.WithLabelValues(TypeIncoming, StatusErr).Inc()
		return err
	}
	duration := time.Since(startTime)

	SummaryResponseTime.Observe(duration.Seconds())
	HistogramResponseTime.Observe(duration.Seconds())

	InFlightTypeMsg.WithLabelValues(TypeIncoming, StatusOk).Inc()

	return nil
}
