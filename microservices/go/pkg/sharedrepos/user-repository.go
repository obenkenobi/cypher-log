package sharedrepos

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmodels"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	baserepos.CRUDRepository[sharedmodels.User, string]
	FindByUserId(ctx context.Context, userId string) (option.Maybe[sharedmodels.User], error)
	FindByAuthId(ctx context.Context, authId string) (option.Maybe[sharedmodels.User], error)
}

type UserRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[sharedmodels.User]
}

func (u UserRepositoryImpl) Create(ctx context.Context, user sharedmodels.User) (sharedmodels.User, error) {
	err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &user)
	return user, err
}

func (u UserRepositoryImpl) Update(ctx context.Context, user sharedmodels.User) (sharedmodels.User, error) {
	err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &user)
	return user, err
}

func (u UserRepositoryImpl) Delete(ctx context.Context, user sharedmodels.User) (sharedmodels.User, error) {
	err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &user)
	return user, err
}

func (u UserRepositoryImpl) FindById(ctx context.Context, id string) (option.Maybe[sharedmodels.User], error) {
	return dshandlers.HandleSingleFind(u.MongoDBHandler, func() (sharedmodels.User, error) {
		user := sharedmodels.User{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.ToChildCtx(ctx), id, &user)
		return user, err
	})
}

func (u UserRepositoryImpl) FindByUserId(ctx context.Context, userId string) (option.Maybe[sharedmodels.User], error) {
	return dshandlers.HandleSingleFind(u.MongoDBHandler, func() (sharedmodels.User, error) {
		user := sharedmodels.User{}
		err := mgm.Coll(u.ModelColl).FirstWithCtx(u.MongoDBHandler.ToChildCtx(ctx), bson.M{"userId": userId}, &user)
		return user, err
	})
}

func (u UserRepositoryImpl) FindByAuthId(ctx context.Context, authId string) (option.Maybe[sharedmodels.User], error) {
	return dshandlers.HandleSingleFind(u.MongoDBHandler, func() (sharedmodels.User, error) {
		user := sharedmodels.User{}
		err := mgm.Coll(u.ModelColl).FirstWithCtx(u.MongoDBHandler.ToChildCtx(ctx), bson.M{"authId": authId}, &user)
		return user, err
	})
}

func NewUserRepositoryImpl(mongoDBHandler *dshandlers.MongoDBHandler) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[sharedmodels.User](sharedmodels.User{}, mongoDBHandler),
	}
}
