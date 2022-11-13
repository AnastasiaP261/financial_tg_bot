package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/financial-tg-bot/metrics"
	"go.uber.org/zap"
)

func (c *Client) SetString(ctx context.Context, key, value string) error {
	stCmd := c.rdb.Set(ctx, key, value, statusesTTL)

	err := stCmd.Err()
	if err != nil {
		logs.Error("set report error", zap.Error(err))
		metrics.InFlightCache.WithLabelValues(metrics.StatusErr).Inc()
		return err
	}
	metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()

	return nil
}

func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	res := c.rdb.Get(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		// здесь все ок так как такого значения просто нет в кеше
		metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()
		return "", nil
	}

	err := res.Err()
	if err != nil {
		logs.Error("set report error", zap.Error(err))
		metrics.InFlightCache.WithLabelValues(metrics.StatusErr).Inc()
		return "", err
	}
	metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()

	return res.Val(), nil
}
