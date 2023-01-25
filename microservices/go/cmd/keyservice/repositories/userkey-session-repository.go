package repositories

import (
	"context"
	"github.com/akrennmair/slice"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
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

func (u UserKeySessionRepositoryImpl) Get(
	ctx context.Context,
	key string,
) (option.Maybe[models.UserKeySession], error) {
	return u.baseRepo.Get(ctx, kvstoreutils.CombineKeySections(u.prefix, key))
}

func (u UserKeySessionRepositoryImpl) Set(
	ctx context.Context,
	key string,
	value models.UserKeySession,
	expiration time.Duration,
) (models.UserKeySession, error) {
	return u.baseRepo.Set(ctx, kvstoreutils.CombineKeySections(u.prefix, key), value, expiration)
}

func (u UserKeySessionRepositoryImpl) Del(ctx context.Context, keys ...string) error {
	nonEmptyKeys := slice.Filter(keys, utils.StringIsNotBlank)
	combinedKeys := slice.Map(nonEmptyKeys, func(key string) string {
		return kvstoreutils.CombineKeySections(u.prefix, key)
	})
	return u.baseRepo.Del(ctx, combinedKeys...)
}

func NewUserKeySessionRepositoryImpl(redisDBHandler *dshandlers.RedisDBHandler) *UserKeySessionRepositoryImpl {
	prefix := kvstoreutils.CombineKeySections(kvStoreKeyPrefix, "userKeySession")
	baseRepo := baserepos.NewKeyValueTimedRepositoryRedis[models.UserKeySession](redisDBHandler)
	return &UserKeySessionRepositoryImpl{prefix: prefix, baseRepo: baseRepo}
}
