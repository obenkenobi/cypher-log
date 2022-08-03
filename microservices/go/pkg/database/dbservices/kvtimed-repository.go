package dbservices

import (
	"context"
	"encoding/json"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"time"
)

type KeyValueTimedRepository[Key any, Value any] interface {
	Get(ctx context.Context, key Key) single.Single[Value]
	GetAsync(ctx context.Context, key Key) single.Single[Value]
	Set(ctx context.Context, key Key, value Value, expiration time.Duration) single.Single[Value]
	SetAsync(ctx context.Context, key Key, expiration time.Duration) single.Single[Value]
}

// KeyValueTimedRepositoryRedis is a Redis implementation of KeyValueTimedRepository
type KeyValueTimedRepositoryRedis[Value any] struct {
	redisDBHandler RedisKeyValueTimedDBHandler
}

func (k KeyValueTimedRepositoryRedis[Value]) Get(ctx context.Context, key string) single.Single[Value] {
	return single.FromSupplier[Value](func() (Value, error) { return k.runGet(ctx, key) })
}

func (k KeyValueTimedRepositoryRedis[Value]) GetAsync(ctx context.Context, key string) single.Single[Value] {
	return single.FromSupplierAsync[Value](func() (Value, error) { return k.runGet(ctx, key) })
}

func (k KeyValueTimedRepositoryRedis[Value]) runGet(ctx context.Context, key string) (Value, error) {
	var val Value
	valJson, err := k.redisDBHandler.GetRedisClient().Get(ctx, key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(valJson), &val)
	}
	return val, err
}

func (k KeyValueTimedRepositoryRedis[Value]) Set(
	ctx context.Context,
	key string,
	value Value,
	expiration time.Duration,
) single.Single[Value] {
	return single.FromSupplier[Value](func() (Value, error) { return k.runSet(ctx, key, value, expiration) })
}

func (k KeyValueTimedRepositoryRedis[Value]) SetAsync(
	ctx context.Context,
	key string,
	value Value,
	expiration time.Duration,
) single.Single[Value] {
	return single.FromSupplierAsync[Value](func() (Value, error) { return k.runSet(ctx, key, value, expiration) })
}

func (k KeyValueTimedRepositoryRedis[Value]) runSet(
	ctx context.Context,
	key string,
	value Value,
	expiration time.Duration,
) (Value, error) {
	valJson, err := json.Marshal(value)
	if err == nil {
		err = k.redisDBHandler.GetRedisClient().Set(ctx, key, string(valJson), expiration).Err()
	}
	return value, err
}

func NewKeyValueTimedRepositoryRedis[Value any](
	redisDBHandler RedisKeyValueTimedDBHandler,
) *KeyValueTimedRepositoryRedis[Value] {
	return &KeyValueTimedRepositoryRedis[Value]{redisDBHandler: redisDBHandler}
}
