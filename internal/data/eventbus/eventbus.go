package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/sqlmerr/astragalaxy/internal/data/redis"
	core_logger "github.com/sqlmerr/astragalaxy/internal/logger"
	"github.com/sqlmerr/astragalaxy/internal/model"
	"go.uber.org/zap"
)

const queueKey = "queue:events"

type envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type rawHandler func(ctx context.Context, payload json.RawMessage) error

type Bus struct {
	rdb      redis.Redis
	log      core_logger.Logger
	handlers map[string][]rawHandler
}

func New(rdb redis.Redis, log core_logger.Logger) *Bus {
	return &Bus{
		rdb:      rdb,
		log:      log,
		handlers: make(map[string][]rawHandler),
	}
}

func (b *Bus) Publish(ctx context.Context, events ...model.Event) {
	for _, ev := range events {
		payload, err := json.Marshal(ev)
		if err != nil {
			b.log.Error("failed to marshal event payload", zap.String("type", ev.EventType()), zap.Error(err))
			continue
		}
		env, err := json.Marshal(envelope{
			Type:    ev.EventType(),
			Payload: payload,
		})
		if err != nil {
			b.log.Error("failed to marshal event envelope", zap.String("type", ev.EventType()), zap.Error(err))
			continue
		}
		if err := b.rdb.LPush(ctx, queueKey, env).Err(); err != nil {
			b.log.Error("failed to push event to redis", zap.String("type", ev.EventType()), zap.Error(err))
		}
	}
}

func Subsribe[T model.Event](b *Bus, handler func(ctx context.Context, ev T) error) { // TODO(update): go 1.27
	var zero T
	eventType := zero.EventType()

	b.handlers[eventType] = append(b.handlers[eventType], func(ctx context.Context, payload json.RawMessage) error {
		var ev T
		if err := json.Unmarshal(payload, &ev); err != nil {
			return err
		}
		return handler(ctx, ev)
	})
}

func (b *Bus) StartWorker(ctx context.Context) {
	for {
		res, err := b.rdb.BRPop(ctx, time.Second*5, queueKey).Result()
		if errors.Is(err, goredis.Nil) {
			continue
		}
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			b.log.Error("redis event bus pop error", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		var env envelope
		if err := json.Unmarshal([]byte(res[1]), &env); err != nil {
			b.log.Error("invalid event envelope", zap.Error(err))
			continue
		}

		for _, h := range b.handlers[env.Type] {
			if err := h(ctx, env.Payload); err != nil {
				b.log.Error("event handler failed", zap.Error(err), zap.String("type", env.Type))
			}
		}
	}
}
