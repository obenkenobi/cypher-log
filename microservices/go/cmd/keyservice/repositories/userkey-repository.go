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

type UserKeyRepository interface {
	baserepos.CRUDRepository[models.UserKey, string]
	FindOneByUserId(ctx context.Context, userId string) single.Single[option.Maybe[models.UserKey]]
	DeleteByUserIdAndGetCount(ctx context.Context, userId string) single.Single[int64]
}

type UserKeyRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[models.UserKey]
}

func (u UserKeyRepositoryImpl) Create(ctx context.Context, model models.UserKey) single.Single[models.UserKey] {
	return single.FromSupplier(func() (models.UserKey, error) {
		err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &model)
		return model, err
	})
}

func (u UserKeyRepositoryImpl) Update(ctx context.Context, model models.UserKey) single.Single[models.UserKey] {
	return single.FromSupplier(func() (models.UserKey, error) {
		err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &model)
		return model, err
	})
}

func (u UserKeyRepositoryImpl) Delete(ctx context.Context, model models.UserKey) single.Single[models.UserKey] {
	return single.FromSupplier(func() (models.UserKey, error) {
		err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &model)
		return model, err
	})
}

func (u UserKeyRepositoryImpl) FindById(ctx context.Context, id string) single.Single[option.Maybe[models.UserKey]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.UserKey, error) {
		model := models.UserKey{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), id, &model)
		return model, err
	})
}

func (u UserKeyRepositoryImpl) FindOneByUserId(
	ctx context.Context,
	userId string,
) single.Single[option.Maybe[models.UserKey]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.UserKey, error) {
		user := models.UserKey{}
		err := mgm.Coll(u.ModelColl).
			FirstWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), bson.M{"userId": userId}, &user)
		return user, err
	})
}

func (u UserKeyRepositoryImpl) DeleteByUserIdAndGetCount(ctx context.Context, userId string) single.Single[int64] {
	return single.FromSupplier(func() (int64, error) {
		res, err := mgm.Coll(u.ModelColl).DeleteMany(ctx, bson.M{"userId": userId})
		deletedCount := option.
			Map(option.Perhaps(res), func(r *mongo.DeleteResult) int64 { return r.DeletedCount }).
			OrElse(-1)
		return deletedCount, err
	})

}

func NewUserKeyRepositoryImpl(mongoDBHandler *dshandlers.MongoDBHandler) *UserKeyRepositoryImpl {
	return &UserKeyRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[models.UserKey](models.UserKey{}, mongoDBHandler),
	}
}

// Delete from userkeys where u.userId = "243234"
