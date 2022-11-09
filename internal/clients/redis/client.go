package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
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
	return res.Err()
}
