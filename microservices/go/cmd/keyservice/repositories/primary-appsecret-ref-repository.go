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

type PrimaryAppSecretRefRepository interface {
	Get(ctx context.Context) single.Single[option.Maybe[models.PrimaryAppSecretRef]]
	Set(ctx context.Context, value models.PrimaryAppSecretRef, expr time.Duration) single.Single[models.PrimaryAppSecretRef]
}

type PrimaryAppSecretRefRepositoryImpl struct {
	prefix   string
	baseRepo baserepos.KeyValueTimedRepository[models.PrimaryAppSecretRef]
}

func (a PrimaryAppSecretRefRepositoryImpl) Get(ctx context.Context) single.Single[option.Maybe[models.PrimaryAppSecretRef]] {
	return a.baseRepo.Get(ctx, a.prefix)
}

func (a PrimaryAppSecretRefRepositoryImpl) Set(
	ctx context.Context,
	value models.PrimaryAppSecretRef,
	expiration time.Duration,
) single.Single[models.PrimaryAppSecretRef] {
	return a.baseRepo.Set(ctx, a.prefix, value, expiration)
}

func NewPrimaryAppSecretRefRepositoryImpl(
	redisDBHandler *dshandlers.RedisDBHandler,
) *PrimaryAppSecretRefRepositoryImpl {
	appSecretKeyPrefix := kvstoreutils.CombineKeySections(kvStoreKeyPrefix, "mainAppSecretRef")
	baseRepo := baserepos.NewKeyValueTimedRepositoryRedis[models.PrimaryAppSecretRef](redisDBHandler)
	return &PrimaryAppSecretRefRepositoryImpl{prefix: appSecretKeyPrefix, baseRepo: baseRepo}
}
