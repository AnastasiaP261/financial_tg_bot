package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/metrics"
	"go.uber.org/zap"
)

var (
	timeout     = time.Second * 15
	statusesTTL = time.Minute * 30
)

type configGetter interface {
	RedisUri() string
	RedisPass() string
}

type Client struct {
	rdb *redis.Client
}

func New(ctx context.Context, config configGetter) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.RedisUri(),
		Password:     config.RedisPass(),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		DB:           0, // use default DB
	})

	return &Client{rdb: rdb}, rdb.Ping(ctx).Err()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	res := c.rdb.Del(ctx, key)

	err := res.Err()
	if err != nil {
		logs.Error("set report error", zap.Error(err))
		metrics.InFlightCache.WithLabelValues(metrics.StatusErr).Inc()
		return err
	}
	metrics.InFlightCache.WithLabelValues(metrics.StatusOk).Inc()

	return nil
}
