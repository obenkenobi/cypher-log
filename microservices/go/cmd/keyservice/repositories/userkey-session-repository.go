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

type UserKeySessionRepository interface {
	baserepos.KeyValueTimedRepository[models.UserKeySession]
}

type UserKeySessionRepositoryImpl struct {
	prefix   string
	baseRepo baserepos.KeyValueTimedRepository[models.UserKeySession]
}

func (a UserKeySessionRepositoryImpl) Get(
	ctx context.Context,
	key string,
) single.Single[option.Maybe[models.UserKeySession]] {
	return a.baseRepo.Get(ctx, kvstoreutils.CombineKeySections(a.prefix, key))
}

func (a UserKeySessionRepositoryImpl) Set(
	ctx context.Context,
	key string,
	value models.UserKeySession,
	expiration time.Duration,
) single.Single[models.UserKeySession] {
	return a.baseRepo.Set(ctx, kvstoreutils.CombineKeySections(a.prefix, key), value, expiration)
}

func NewUserKeySessionRepositoryImpl(redisDBHandler *dshandlers.RedisDBHandler) *UserKeySessionRepositoryImpl {
	prefix := kvstoreutils.CombineKeySections(kvStoreKeyPrefix, "userKeySession")
	baseRepo := baserepos.NewKeyValueTimedRepositoryRedis[models.UserKeySession](redisDBHandler)
	return &UserKeySessionRepositoryImpl{prefix: prefix, baseRepo: baseRepo}
}
