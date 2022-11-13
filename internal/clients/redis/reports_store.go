package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/report"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/financial-tg-bot/metrics"
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

func (c *Client) SetReport(ctx context.Context, key string, value report.Report) error {
	items := make([]ReportItem, len(value.Items))
	for i := range value.Items {
		items[i] = ReportItem(value.Items[i])
	}

	r := Report{
		Items:    items,
		FromDate: value.FromDate,
	}

	var nullDur time.Duration
	stCmd := c.rdb.Set(ctx, key, r, nullDur)

	err := stCmd.Err()
	if err != nil {
		logs.Error("set r error", zap.Error(err))
		metrics.InFlightCache.WithLabelValues(metrics.StatusErr).Inc()
		return err
	}
	metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()

	return nil
}

func (c *Client) GetReport(ctx context.Context, key string) (report.Report, error) {
	res := c.rdb.Get(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		// здесь все ок так как такого значения просто нет в кеше
		metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()
		return report.Report{}, nil
	}

	var r Report
	if err := res.Scan(&r); err != nil {
		logs.Error("set r error", zap.Error(err))
		metrics.InFlightCache.WithLabelValues(metrics.StatusErr).Inc()
		return report.Report{}, errors.Wrap(err, "scanning err")
	}
	metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()

	items := make([]report.ReportItem, len(r.Items))
	for i := range r.Items {
		items[i] = report.ReportItem(r.Items[i])
	}

	return report.Report{
		Items:    items,
		FromDate: r.FromDate,
	}, nil
}
