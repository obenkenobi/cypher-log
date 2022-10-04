package repositories

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/kvstoreutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"time"
)

type AppSecretRepository interface {
	baserepos.KeyValueTimedRepository[models.AppSecret]
}

type AppSecretRepositoryImpl struct {
	prefix   string
	baseRepo baserepos.KeyValueTimedRepository[models.AppSecret]
}

func (a AppSecretRepositoryImpl) Get(ctx context.Context, key string) single.Single[option.Maybe[models.AppSecret]] {
	return a.baseRepo.Get(ctx, kvstoreutils.CombineKeySections(a.prefix, key))
}

func (a AppSecretRepositoryImpl) Set(
	ctx context.Context,
	key string,
	value models.AppSecret,
	expiration time.Duration,
) single.Single[models.AppSecret] {
	return a.baseRepo.Set(ctx, kvstoreutils.CombineKeySections(a.prefix, key), value, expiration)
}

func NewAppSecretRepositoryImpl(redisDBHandler *dshandlers.RedisKeyValueTimedDBHandler) *AppSecretRepositoryImpl {
	prefix := kvstoreutils.CombineKeySections(kvStoreKeyPrefix, "appSecret")
	baseRepo := baserepos.NewKeyValueTimedRepositoryRedis[models.AppSecret](redisDBHandler)
	return &AppSecretRepositoryImpl{prefix: prefix, baseRepo: baseRepo}
}
