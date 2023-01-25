package repositories

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/kvstoreutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"time"
)

type AccessTokenHolderRepository interface {
	baserepos.KeyValueTimedRepository[models.AccessTokenHolder]
}

type AccessTokenHolderRepositoryImpl struct {
	prefix   string
	baseRepo baserepos.KeyValueTimedRepository[models.AccessTokenHolder]
}

func (a AccessTokenHolderRepositoryImpl) Get(
	ctx context.Context,
	key string,
) (option.Maybe[models.AccessTokenHolder], error) {
	return a.baseRepo.Get(ctx, kvstoreutils.CombineKeySections(a.prefix, key))
}

func (a AccessTokenHolderRepositoryImpl) Set(
	ctx context.Context,
	key string,
	value models.AccessTokenHolder,
	expiration time.Duration,
) (models.AccessTokenHolder, error) {
	return a.baseRepo.Set(ctx, kvstoreutils.CombineKeySections(a.prefix, key), value, expiration)
}

func NewAccessTokenHolderRepositoryImpl(redisDBHandler *dshandlers.RedisDBHandler) *AccessTokenHolderRepositoryImpl {
	prefix := kvstoreutils.CombineKeySections(kvStoreKeyPrefix, "accessTokenHolder")
	baseRepo := baserepos.NewKeyValueTimedRepositoryRedis[models.AccessTokenHolder](redisDBHandler)
	return &AccessTokenHolderRepositoryImpl{prefix: prefix, baseRepo: baseRepo}
}
