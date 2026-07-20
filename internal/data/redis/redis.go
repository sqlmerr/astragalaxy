package redis

import (
	"context"
	"time"
)

type Redis interface {
	Get(ctx context.Context, key string) StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) StringCmd
	Del(ctx context.Context, keys ...string) IntCmd
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

type IntCmd interface {
	Cmd
	Val() int64
	Result() (int64, error)
}
