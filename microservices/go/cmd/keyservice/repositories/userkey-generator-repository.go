package repositories

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
)

type UserKeyGeneratorRepository interface {
	baserepos.CRUDRepository[models.UserKeyGenerator, string]
	FindOneByUserId(ctx context.Context, userId string) (option.Maybe[models.UserKeyGenerator], error)
	DeleteByUserIdAndGetCount(ctx context.Context, userId string) (int64, error)
}

type UserKeyGeneratorRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[models.UserKeyGenerator]
}

func (u UserKeyGeneratorRepositoryImpl) Create(
	ctx context.Context,
	model models.UserKeyGenerator,
) (models.UserKeyGenerator, error) {
	err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
	return model, err
}

func (u UserKeyGeneratorRepositoryImpl) Update(
	ctx context.Context,
	model models.UserKeyGenerator,
) (models.UserKeyGenerator, error) {
	err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
	return model, err
}

func (u UserKeyGeneratorRepositoryImpl) Delete(
	ctx context.Context,
	model models.UserKeyGenerator,
) (models.UserKeyGenerator, error) {
	err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
	return model, err
}

func (u UserKeyGeneratorRepositoryImpl) FindById(
	ctx context.Context,
	id string,
) (option.Maybe[models.UserKeyGenerator], error) {
	return dshandlers.HandleSingleFind(u.MongoDBHandler, func() (models.UserKeyGenerator, error) {
		model := models.UserKeyGenerator{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.ToChildCtx(ctx), id, &model)
		return model, err
	})
}

func (u UserKeyGeneratorRepositoryImpl) FindOneByUserId(
	ctx context.Context,
	userId string,
) (option.Maybe[models.UserKeyGenerator], error) {
	return dshandlers.HandleSingleFind(u.MongoDBHandler, func() (models.UserKeyGenerator, error) {
		user := models.UserKeyGenerator{}
		err := mgm.Coll(u.ModelColl).
			FirstWithCtx(u.MongoDBHandler.ToChildCtx(ctx), bson.M{"userId": userId}, &user)
		return user, err
	})
}

func (u UserKeyGeneratorRepositoryImpl) DeleteByUserIdAndGetCount(
	ctx context.Context,
	userId string,
) (int64, error) {
	res, err := mgm.Coll(u.ModelColl).DeleteMany(u.MongoDBHandler.ToChildCtx(ctx), bson.M{"userId": userId})
	if res != nil {
		return res.DeletedCount, err
	}
	return -1, err

}

func NewUserKeyRepositoryImpl(mongoDBHandler *dshandlers.MongoDBHandler) *UserKeyGeneratorRepositoryImpl {
	return &UserKeyGeneratorRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[models.UserKeyGenerator](
			models.UserKeyGenerator{},
			mongoDBHandler,
		),
	}
}
