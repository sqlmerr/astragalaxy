package redis_goredis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sqlmerr/astragalaxy/internal/data/redis"
)

type Client struct {
	*goredis.Client
}

func NewClient(ctx context.Context, config Config) (*Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Client{rdb}, nil
}

func (c *Client) Get(ctx context.Context, key string) redis.StringCmd {
	cmd := c.Client.Get(ctx, key)
	return goredisStringCmd{cmd}
}

func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) redis.StringCmd {
	cmd := c.Client.Set(ctx, key, value, expiration)
	return goredisStatusCmd{cmd}
}

func (c *Client) Del(ctx context.Context, keys ...string) redis.IntCmd {
	cmd := c.Client.Del(ctx, keys...)
	return goredisIntCmd{cmd}
}
