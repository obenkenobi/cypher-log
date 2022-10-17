package repositories

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserKeyGeneratorRepository interface {
	baserepos.CRUDRepository[models.UserKeyGenerator, string]
	FindOneByUserId(ctx context.Context, userId string) single.Single[option.Maybe[models.UserKeyGenerator]]
	DeleteByUserIdAndGetCount(ctx context.Context, userId string) single.Single[int64]
}

type UserKeyGeneratorRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[models.UserKeyGenerator]
}

func (u UserKeyGeneratorRepositoryImpl) Create(
	ctx context.Context,
	model models.UserKeyGenerator,
) single.Single[models.UserKeyGenerator] {
	return single.FromSupplierCached(func() (models.UserKeyGenerator, error) {
		err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
		return model, err
	})
}

func (u UserKeyGeneratorRepositoryImpl) Update(
	ctx context.Context,
	model models.UserKeyGenerator,
) single.Single[models.UserKeyGenerator] {
	return single.FromSupplierCached(func() (models.UserKeyGenerator, error) {
		err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
		return model, err
	})
}

func (u UserKeyGeneratorRepositoryImpl) Delete(
	ctx context.Context,
	model models.UserKeyGenerator,
) single.Single[models.UserKeyGenerator] {
	return single.FromSupplierCached(func() (models.UserKeyGenerator, error) {
		err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
		return model, err
	})
}

func (u UserKeyGeneratorRepositoryImpl) FindById(
	ctx context.Context,
	id string,
) single.Single[option.Maybe[models.UserKeyGenerator]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.UserKeyGenerator, error) {
		model := models.UserKeyGenerator{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.ToChildCtx(ctx), id, &model)
		return model, err
	})
}

func (u UserKeyGeneratorRepositoryImpl) FindOneByUserId(
	ctx context.Context,
	userId string,
) single.Single[option.Maybe[models.UserKeyGenerator]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.UserKeyGenerator, error) {
		user := models.UserKeyGenerator{}
		err := mgm.Coll(u.ModelColl).
			FirstWithCtx(u.MongoDBHandler.ToChildCtx(ctx), bson.M{"userId": userId}, &user)
		return user, err
	})
}

func (u UserKeyGeneratorRepositoryImpl) DeleteByUserIdAndGetCount(
	ctx context.Context,
	userId string,
) single.Single[int64] {
	return single.FromSupplierCached(func() (int64, error) {
		res, err := mgm.Coll(u.ModelColl).DeleteMany(u.MongoDBHandler.ToChildCtx(ctx), bson.M{"userId": userId})
		deletedCount := option.
			Map(option.Perhaps(res), func(r *mongo.DeleteResult) int64 { return r.DeletedCount }).
			OrElse(-1)
		return deletedCount, err
	})

}

func NewUserKeyRepositoryImpl(mongoDBHandler *dshandlers.MongoDBHandler) *UserKeyGeneratorRepositoryImpl {
	return &UserKeyGeneratorRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[models.UserKeyGenerator](
			models.UserKeyGenerator{},
			mongoDBHandler,
		),
	}
}
