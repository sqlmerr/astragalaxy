package testutil

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sqlmerr/astragalaxy/internal/data/redis"
	redis_goredis "github.com/sqlmerr/astragalaxy/internal/data/redis/goredis"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

func SetupRedis(ctx context.Context) (redis.Redis, func(), error) {
	redisContainer, err := tcredis.Run(ctx, "redis:8-alpine")
	if err != nil {
		return nil, nil, fmt.Errorf("start redis container: %w", err)
	}

	connStr, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		_ = testcontainers.TerminateContainer(redisContainer)
		return nil, nil, fmt.Errorf("get redis connection string: %w", err)
	}

	opts, err := goredis.ParseURL(connStr)
	if err != nil {
		_ = testcontainers.TerminateContainer(redisContainer)
		return nil, nil, fmt.Errorf("parse redis url: %w", err)
	}

	rdb := goredis.NewClient(opts)

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		_ = testcontainers.TerminateContainer(redisContainer)
		return nil, nil, fmt.Errorf("ping redis: %w", err)
	}

	cleanup := func() {
		_ = rdb.Close()
		_ = testcontainers.TerminateContainer(redisContainer)
	}

	return &redis_goredis.Client{Client: rdb}, cleanup, nil
}
