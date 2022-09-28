package dshandlers

import (
	"github.com/go-redis/redis/v9"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
)

type KeyValueTimedDSHandler interface {
	DataSourceHandler
}

type RedisKeyValueTimedDBHandler struct {
	redisClient *redis.Client
}

func (r RedisKeyValueTimedDBHandler) GetRedisClient() *redis.Client {
	return r.redisClient
}

func (r RedisKeyValueTimedDBHandler) IsNotFoundError(err error) bool {
	return err == redis.Nil
}

func NewRedisKeyValueTimedDBHandler(redisConf conf.RedisConf) *RedisKeyValueTimedDBHandler {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConf.GetAddress(),
		Password: redisConf.GetPassword(),
		DB:       redisConf.GetDB(),
	})
	return &RedisKeyValueTimedDBHandler{redisClient: rdb}
}
