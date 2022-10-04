package dshandlers

import (
	"github.com/go-redis/redis/v9"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
)

type KeyValueTimedDSHandler interface {
	DataSourceHandler
}

type RedisDBHandler struct {
	redisClient *redis.Client
}

func (r RedisDBHandler) GetRedisClient() *redis.Client {
	return r.redisClient
}

func (r RedisDBHandler) IsNotFoundError(err error) bool {
	return err == redis.Nil
}

func NewRedisDBHandler(redisConf conf.RedisConf) *RedisDBHandler {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConf.GetAddress(),
		Password: redisConf.GetPassword(),
		DB:       redisConf.GetDB(),
	})
	return &RedisDBHandler{redisClient: rdb}
}
