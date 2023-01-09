package repositories

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/mgmtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	baserepos.CRUDRepository[models.User, string]
	FindByAuthIdAndNotToBeDeleted(ctx context.Context, authId string) (option.Maybe[models.User], error)
	FindByUsernameAndNotToBeDeleted(ctx context.Context, username string) (option.Maybe[models.User], error)
	SampleUndistributedUsers(ctx context.Context, size int64) ([]models.User, error)
}

type UserRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[models.User]
}

func (u UserRepositoryImpl) Create(ctx context.Context, model models.User) (models.User, error) {
	err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
	return model, err
}

func (u UserRepositoryImpl) Update(ctx context.Context, model models.User) (models.User, error) {
	err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
	return model, err
}

func (u UserRepositoryImpl) Delete(ctx context.Context, model models.User) (models.User, error) {
	err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
	return model, err
}

func (u UserRepositoryImpl) FindById(ctx context.Context, id string) (option.Maybe[models.User], error) {
	return dshandlers.HandleSingleFind(u.MongoDBHandler, func() (models.User, error) {
		model := models.User{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.ToChildCtx(ctx), id, &model)
		return model, err
	})
}

func (u UserRepositoryImpl) FindByAuthIdAndNotToBeDeleted(
	ctx context.Context,
	authId string,
) (option.Maybe[models.User], error) {
	return dshandlers.HandleSingleFind(u.MongoDBHandler, func() (models.User, error) {
		user := models.User{}
		err := mgm.Coll(u.ModelColl).FirstWithCtx(u.MongoDBHandler.ToChildCtx(ctx),
			bson.M{
				"authId":      authId,
				"toBeDeleted": bson.M{operator.Ne: true}},
			&user)
		return user, err
	})
}

func (u UserRepositoryImpl) FindByUsernameAndNotToBeDeleted(
	ctx context.Context,
	username string,
) (option.Maybe[models.User], error) {
	return dshandlers.HandleSingleFind(u.MongoDBHandler, func() (models.User, error) {
		user := models.User{}
		err := mgm.Coll(u.ModelColl).FirstWithCtx(
			u.MongoDBHandler.ToChildCtx(ctx),
			bson.M{
				"userName":    username,
				"toBeDeleted": bson.M{operator.Ne: true},
			},
			&user,
		)
		return user, err
	})
}

func (u UserRepositoryImpl) SampleUndistributedUsers(ctx context.Context, size int64) ([]models.User, error) {
	childCtx := u.MongoDBHandler.ToChildCtx(ctx)
	cursor, err := mgm.Coll(u.ModelColl).Aggregate(ctx, mongo.Pipeline{
		{{operator.Match, bson.D{{"distributed", bson.D{{operator.Ne, true}}}}}},
		{{operator.Sample, bson.D{{"size", size}}}},
	})
	return mgmtools.HandleFindManyRes[models.User](childCtx, cursor, err)
}

func NewUserRepositoryImpl(mongoDBHandler *dshandlers.MongoDBHandler) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[models.User](models.User{}, mongoDBHandler),
	}
}
