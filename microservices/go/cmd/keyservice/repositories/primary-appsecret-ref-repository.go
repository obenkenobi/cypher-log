package repositories

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/kvstoreutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"time"
)

type PrimaryAppSecretRefRepository interface {
	Get(ctx context.Context) (option.Maybe[models.PrimaryAppSecretRef], error)
	Set(ctx context.Context, value models.PrimaryAppSecretRef, expr time.Duration) (models.PrimaryAppSecretRef, error)
}

type PrimaryAppSecretRefRepositoryImpl struct {
	key      string
	baseRepo baserepos.KeyValueTimedRepository[models.PrimaryAppSecretRef]
}

func (a PrimaryAppSecretRefRepositoryImpl) Get(ctx context.Context) (option.Maybe[models.PrimaryAppSecretRef], error) {
	return a.baseRepo.Get(ctx, a.key)
}

func (a PrimaryAppSecretRefRepositoryImpl) Set(
	ctx context.Context,
	value models.PrimaryAppSecretRef,
	expiration time.Duration,
) (models.PrimaryAppSecretRef, error) {
	return a.baseRepo.Set(ctx, a.key, value, expiration)
}

func NewPrimaryAppSecretRefRepositoryImpl(
	redisDBHandler *dshandlers.RedisDBHandler,
) *PrimaryAppSecretRefRepositoryImpl {
	key := kvstoreutils.CombineKeySections(kvStoreKeyPrefix, "mainAppSecretRef")
	baseRepo := baserepos.NewKeyValueTimedRepositoryRedis[models.PrimaryAppSecretRef](redisDBHandler)
	return &PrimaryAppSecretRefRepositoryImpl{key: key, baseRepo: baseRepo}
}
