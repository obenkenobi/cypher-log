package baserepos

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"time"
)

type KeyValueTimedRepository[Value any] interface {
	Get(ctx context.Context, key string) (option.Maybe[Value], error)
	Set(ctx context.Context, key string, value Value, expiration time.Duration) (Value, error)
}

// KeyValueTimedRepositoryRedis is a Redis implementation of KeyValueTimedRepository
type KeyValueTimedRepositoryRedis[Value any] struct {
	redisDBHandler *dshandlers.RedisDBHandler
}

func (k KeyValueTimedRepositoryRedis[Value]) Get(ctx context.Context, key string) (option.Maybe[Value], error) {
	if utils.StringIsBlank(key) {
		return option.None[Value](), errors.New("key cannot be empty")
	}
	valJson, err := k.redisDBHandler.GetRedisClient().Get(ctx, key).Result()
	if k.redisDBHandler.IsNotFoundError(err) {
		return option.None[Value](), nil
	} else if err != nil {
		return option.None[Value](), err
	} else {
		var val Value
		err = json.Unmarshal([]byte(valJson), &val)
		return option.Perhaps(val), err
	}
}

func (k KeyValueTimedRepositoryRedis[Value]) Set(
	ctx context.Context,
	key string,
	value Value,
	expiration time.Duration,
) (Value, error) {
	if utils.StringIsBlank(key) {
		return value, errors.New("key cannot be empty")
	}
	valJson, err := json.Marshal(value)
	if err == nil {
		err = k.redisDBHandler.GetRedisClient().Set(ctx, key, string(valJson), expiration).Err()
	}
	return value, err
}

func NewKeyValueTimedRepositoryRedis[Value any](
	redisDBHandler *dshandlers.RedisDBHandler,
) KeyValueTimedRepository[Value] {
	return &KeyValueTimedRepositoryRedis[Value]{redisDBHandler: redisDBHandler}
}
