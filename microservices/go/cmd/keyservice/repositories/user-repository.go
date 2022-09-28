package repositories

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/database/dbservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	dbservices.CRUDRepository[models.User, string]
	FindByUserId(ctx context.Context, userId string) single.Single[option.Maybe[models.User]]
}

type UserRepositoryImpl struct {
	dbservices.BaseRepositoryMongo[models.User]
}

func (u UserRepositoryImpl) Create(ctx context.Context, user models.User) single.Single[models.User] {
	return single.FromSupplier(func() (models.User, error) {
		err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &user)
		return user, err
	})
}

func (u UserRepositoryImpl) Update(ctx context.Context, user models.User) single.Single[models.User] {
	return single.FromSupplier(func() (models.User, error) {
		err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &user)
		return user, err
	})
}

func (u UserRepositoryImpl) Delete(ctx context.Context, user models.User) single.Single[models.User] {
	return single.FromSupplier(func() (models.User, error) {
		err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &user)
		return user, err
	})
}

func (u UserRepositoryImpl) FindById(ctx context.Context, id string) single.Single[option.Maybe[models.User]] {
	return dbservices.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.User, error) {
		user := models.User{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), id, &user)
		return user, err
	})
}

func (u UserRepositoryImpl) FindByUserId(ctx context.Context, userId string) single.Single[option.Maybe[models.User]] {
	return dbservices.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.User, error) {
		user := models.User{}
		err := mgm.Coll(u.ModelColl).FirstWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), bson.M{"userId": userId}, &user)
		return user, err
	})
}

func NewUserMongoRepository(mongoDBHandler *dbservices.MongoDBHandler) UserRepository {
	return &UserRepositoryImpl{
		BaseRepositoryMongo: *dbservices.NewBaseRepositoryMongo[models.User](models.User{}, mongoDBHandler),
	}
}
