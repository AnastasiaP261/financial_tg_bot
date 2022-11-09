package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func (c *Client) SetString(ctx context.Context, key, value string) error {
	stCmd := c.rdb.Set(ctx, key, value, statusesTTL)
	return stCmd.Err()
}

func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	res := c.rdb.Get(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		return "", nil
	}

	return res.Val(), res.Err()
}
