package redis

import (
	"context"
	"time"
)

type Redis interface {
	Get(ctx context.Context, key string) StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) StringCmd
	Del(ctx context.Context, keys ...string) IntCmd
	LPush(ctx context.Context, key string, values ...any) IntCmd
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) StringSliceCmd

	Close() error
}

type Cmd interface {
	Args() []any
	String() string
	Err() error
}

type StringCmd interface {
	Cmd
	Result() (string, error)
	Val() string
}

type StringSliceCmd interface {
	Cmd
	Result() ([]string, error)
	Val() []string
}

type IntCmd interface {
	Cmd
	Val() int64
	Result() (int64, error)
}
