package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
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

func (r Report) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Report) UnmarshalBinary(data []byte) error {
	var v Report
	err := json.Unmarshal(data, &v)
	if err != nil {
		return errors.Wrap(err, "unmarshalling error")
	}

	r.FromDate = v.FromDate
	r.Items = v.Items

	return nil
}

func (c *Client) SetReport(ctx context.Context, key string, value purchases.Report) error {
	items := make([]ReportItem, len(value.Items))
	for i := range value.Items {
		items[i] = ReportItem(value.Items[i])
	}

	report := Report{
		Items:    items,
		FromDate: value.FromDate,
	}

	var nullDur time.Duration
	stCmd := c.rdb.Set(ctx, key, report, nullDur)

	err := stCmd.Err()
	if err != nil {
		logs.Error("set report error", zap.Error(err))
		metrics.InFlightCache.WithLabelValues(metrics.StatusErr).Inc()
		return err
	}
	metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()

	return nil
}

func (c *Client) GetReport(ctx context.Context, key string) (purchases.Report, error) {
	res := c.rdb.Get(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		// здесь все ок так как такого значения просто нет в кеше
		metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()
		return purchases.Report{}, nil
	}

	var report Report
	if err := res.Scan(&report); err != nil {
		logs.Error("set report error", zap.Error(err))
		metrics.InFlightCache.WithLabelValues(metrics.StatusErr).Inc()
		return purchases.Report{}, errors.Wrap(err, "scanning err")
	}
	metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()

	items := make([]purchases.ReportItem, len(report.Items))
	for i := range report.Items {
		items[i] = purchases.ReportItem(report.Items[i])
	}

	return purchases.Report{
		Items:    items,
		FromDate: report.FromDate,
	}, nil
}
