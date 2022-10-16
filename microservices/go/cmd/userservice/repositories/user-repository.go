package repositories

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/operator"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	baserepos.CRUDRepository[models.User, string]
	FindByAuthIdAndNotToBeDeleted(ctx context.Context, authId string) single.Single[option.Maybe[models.User]]
	FindByUsernameAndNotToBeDeleted(ctx context.Context, username string) single.Single[option.Maybe[models.User]]
	SampleUndistributedUsers(ctx context.Context, size int64) stream.Observable[models.User]
}

type UserRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[models.User]
}

func (u UserRepositoryImpl) Create(ctx context.Context, model models.User) single.Single[models.User] {
	return single.FromSupplierCached(func() (models.User, error) {
		err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &model)
		return model, err
	})
}

func (u UserRepositoryImpl) Update(ctx context.Context, model models.User) single.Single[models.User] {
	return single.FromSupplierCached(func() (models.User, error) {
		err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &model)
		return model, err
	})
}

func (u UserRepositoryImpl) Delete(ctx context.Context, model models.User) single.Single[models.User] {
	return single.FromSupplierCached(func() (models.User, error) {
		err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &model)
		return model, err
	})
}

func (u UserRepositoryImpl) FindById(ctx context.Context, id string) single.Single[option.Maybe[models.User]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.User, error) {
		model := models.User{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), id, &model)
		return model, err
	})
}

func (u UserRepositoryImpl) FindByAuthIdAndNotToBeDeleted(
	ctx context.Context,
	authId string,
) single.Single[option.Maybe[models.User]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.User, error) {
		user := models.User{}
		err := mgm.Coll(u.ModelColl).FirstWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx),
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
) single.Single[option.Maybe[models.User]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.User, error) {
		user := models.User{}
		err := mgm.Coll(u.ModelColl).FirstWithCtx(
			u.MongoDBHandler.GetChildDBCtx(ctx),
			bson.M{
				"userName":    username,
				"toBeDeleted": bson.M{operator.Ne: true},
			},
			&user,
		)
		return user, err
	})
}

func (u UserRepositoryImpl) SampleUndistributedUsers(ctx context.Context, size int64) stream.Observable[models.User] {
	return stream.FlatMap(stream.Just(any(true)), func(_ any) stream.Observable[models.User] {
		var results []models.User
		cursor, err := mgm.Coll(u.ModelColl).Aggregate(ctx, mongo.Pipeline{
			{{operator.Match, bson.D{{"distributed", bson.D{{operator.Ne, true}}}}}},
			{{operator.Sample, bson.D{{"size", size}}}},
		})
		if err != nil {
			return stream.Error[models.User](err)
		}
		if err = cursor.All(context.TODO(), &results); err != nil {
			return stream.Error[models.User](err)
		}
		return stream.FromSlice(results)
	})
}

func NewUserRepositoryImpl(mongoDBHandler *dshandlers.MongoDBHandler) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[models.User](models.User{}, mongoDBHandler),
	}
}
