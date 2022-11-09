package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/metrics"
	"go.uber.org/zap"
)

type ReportItem struct {
	PurchaseCategory string  `json:"purchaseCategory"`
	Summa            float64 `json:"summa"`
}

type Report struct {
	Items    []ReportItem `json:"items"`
	FromDate time.Time    `json:"fromDate"` // дата начала выборки данных в отчете
}

func (c *Client) SetReport(ctx context.Context, key string, value Report) error {
	var nullDur time.Duration
	stCmd := c.rdb.Set(ctx, key, value, nullDur)

	err := stCmd.Err()
	if err != nil {
		logs.Error("set report error", zap.Error(err))
		metrics.InFlightCache.WithLabelValues(metrics.StatusErr).Inc()
		return err
	}
	metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()

	return nil
}

func (c *Client) GetReport(ctx context.Context, key string) (Report, error) {
	res := c.rdb.Get(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		// здесь все ок так как такого значения просто нет в кеше
		metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()
		return Report{}, nil
	}

	var report Report
	if err := res.Scan(&report); err != nil {
		logs.Error("set report error", zap.Error(err))
		metrics.InFlightCache.WithLabelValues(metrics.StatusErr).Inc()
		return Report{}, errors.Wrap(err, "scanning err")
	}
	metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()

	return report, res.Err()
}
