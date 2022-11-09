package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

type ReportItem struct {
	PurchaseCategory string  `json:"purchaseCategory"`
	Summa            float64 `json:"summa"`
}

type FromDate struct {
	Y int `json:"y"`
	M int `json:"m"`
	D int `json:"d"`
}

type Report struct {
	Items    []ReportItem `json:"items"`
	CreateTo time.Time    `json:"createTo"` // дата создания отчета
	FromDate FromDate     `json:"fromDate"` // дата начала выборки данных в отчете
}

func (c *Client) SetReport(ctx context.Context, key string, value Report) error {
	var nullDur time.Duration
	stCmd := c.rdb.Set(ctx, key, value, nullDur)

	return stCmd.Err()
}

func (c *Client) GetReport(ctx context.Context, key string) (Report, error) {
	res := c.rdb.Get(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		return Report{}, nil
	}

	var report Report
	if err := res.Scan(&report); err != nil {
		return Report{}, errors.Wrap(err, "scanning err")
	}

	return report, res.Err()
}
