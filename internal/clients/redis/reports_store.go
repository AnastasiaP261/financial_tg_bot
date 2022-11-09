package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

type ReportItem struct {
	PurchaseCategory string
	Summa            float64
}

type Report struct {
	Items []ReportItem
	Date  time.Time // дата создания отчета
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
