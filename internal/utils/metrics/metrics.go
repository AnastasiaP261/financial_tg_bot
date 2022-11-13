package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
