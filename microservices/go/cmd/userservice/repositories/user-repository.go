package repositories

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	baserepos.CRUDRepository[models.User, string]
	FindByAuthId(ctx context.Context, authId string) single.Single[option.Maybe[models.User]]
	FindByUsername(ctx context.Context, username string) single.Single[option.Maybe[models.User]]
}

type UserRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[models.User]
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
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.User, error) {
		user := models.User{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), id, &user)
		return user, err
	})
}

func (u UserRepositoryImpl) FindByAuthId(ctx context.Context, authId string) single.Single[option.Maybe[models.User]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.User, error) {
		user := models.User{}
		err := mgm.Coll(u.ModelColl).FirstWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), bson.M{"authId": authId}, &user)
		return user, err
	})
}

func (u UserRepositoryImpl) FindByUsername(
	ctx context.Context,
	username string,
) single.Single[option.Maybe[models.User]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.User, error) {
		user := models.User{}
		err := mgm.Coll(u.ModelColl).FirstWithCtx(
			u.MongoDBHandler.GetChildDBCtx(ctx),
			bson.M{"userName": username},
			&user,
		)
		return user, err
	})
}

func NewUserMongoRepository(mongoDBHandler *dshandlers.MongoDBHandler) UserRepository {
	return &UserRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[models.User](models.User{}, mongoDBHandler),
	}
}
