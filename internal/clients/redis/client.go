package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	timeout = time.Second * 15
	ttl     = time.Minute * 30
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

func (c *Client) Set(ctx context.Context, key, value string) error {
	stCmd := c.rdb.Set(ctx, key, value, ttl)
	return stCmd.Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	res := c.rdb.Get(ctx, key)
	return res.String(), res.Err()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	res := c.rdb.Del(ctx, key)
	return res.Err()
}
