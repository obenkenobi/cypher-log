package repositories

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/database/dbservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	dbservices.CRUDRepository[*models.User, string]
	FindByAuthId(ctx context.Context, authId string) single.Single[option.Maybe[models.User]]
	FindByAuthIdAsync(ctx context.Context, authId string) single.Single[option.Maybe[models.User]]
	FindByUsername(ctx context.Context, username string) single.Single[option.Maybe[models.User]]
	FindByUsernameAsync(ctx context.Context, username string) single.Single[option.Maybe[models.User]]
}

type UserRepositoryImpl struct {
	dbservices.CRUDRepositoryMongo[*models.User, string]
}

func (u UserRepositoryImpl) FindByAuthId(ctx context.Context, authId string) single.Single[option.Maybe[models.User]] {
	return dbservices.ObserveOptionalSingleQuery(u.MongoDBHandler, func() (models.User, error) {
		return u.runFindByAuthId(ctx, authId)
	})
}

func (u UserRepositoryImpl) FindByAuthIdAsync(
	ctx context.Context,
	authId string,
) single.Single[option.Maybe[models.User]] {
	return dbservices.ObserveOptionalSingleQueryAsync(u.MongoDBHandler, func() (models.User, error) {
		return u.runFindByAuthId(ctx, authId)
	})
}

func (u UserRepositoryImpl) runFindByAuthId(ctx context.Context, authId string) (models.User, error) {
	user := models.User{}
	err := mgm.Coll(u.ModelColumn).FirstWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), bson.M{"authId": authId}, &user)
	return user, err
}

func (u UserRepositoryImpl) FindByUsername(
	ctx context.Context,
	username string,
) single.Single[option.Maybe[models.User]] {
	return dbservices.ObserveOptionalSingleQuery(u.MongoDBHandler, func() (models.User, error) {
		return u.runFindByUsername(ctx, username)
	})
}

func (u UserRepositoryImpl) FindByUsernameAsync(
	ctx context.Context,
	username string,
) single.Single[option.Maybe[models.User]] {
	return dbservices.ObserveOptionalSingleQueryAsync(u.MongoDBHandler, func() (models.User, error) {
		return u.runFindByUsername(ctx, username)
	})
}

func (u UserRepositoryImpl) runFindByUsername(ctx context.Context, username string) (models.User, error) {
	user := models.User{}
	err := mgm.Coll(u.ModelColumn).FirstWithCtx(
		u.MongoDBHandler.GetChildDBCtx(ctx),
		bson.M{"userName": username},
		&user,
	)
	return user, err
}

func NewUserMongoRepository(mongoDBHandler *dbservices.MongoDBHandler) UserRepository {
	return &UserRepositoryImpl{
		CRUDRepositoryMongo: *dbservices.NewRepositoryMongoImpl[*models.User, string](
			&models.User{},
			mongoDBHandler,
		),
	}
}
