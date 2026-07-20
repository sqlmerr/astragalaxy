package redis_goredis

import (
	goredis "github.com/redis/go-redis/v9"
)

type goredisStringCmd struct {
	*goredis.StringCmd
}

type goredisStatusCmd struct {
	*goredis.StatusCmd
}

type goredisIntCmd struct {
	*goredis.IntCmd
}
