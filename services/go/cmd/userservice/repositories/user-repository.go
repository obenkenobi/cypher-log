package repositories

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	database.Repository[*models.User, string]
	FindByAuthId(ctx context.Context, authId string) stream.Observable[option.Maybe[*models.User]]
	FindByUsername(ctx context.Context, username string) stream.Observable[option.Maybe[*models.User]]
}

type UserRepositoryImpl struct {
	database.RepositoryMongoImpl[*models.User, string]
}

func (u UserRepositoryImpl) FindByAuthId(
	ctx context.Context,
	authId string,
) stream.Observable[option.Maybe[*models.User]] {
	return database.ObserveOptionalSingleQuery(u.MongoDBHandler, func() (*models.User, error) {
		user := &models.User{}
		err := mgm.Coll(u.ModelColumn).FirstWithCtx(ctx, bson.M{"authId": authId}, user)
		return user, err
	})
}

func (u UserRepositoryImpl) FindByUsername(
	ctx context.Context,
	username string,
) stream.Observable[option.Maybe[*models.User]] {
	return database.ObserveOptionalSingleQuery(u.MongoDBHandler, func() (*models.User, error) {
		user := &models.User{}
		err := mgm.Coll(u.ModelColumn).FirstWithCtx(ctx, bson.M{"userName": username}, user)
		return user, err
	})
}

func NewUserMongoRepository(mongoDBHandler *database.MongoDBHandler) UserRepository {
	return &UserRepositoryImpl{
		RepositoryMongoImpl: *database.NewRepositoryMongoImpl[*models.User, string](&models.User{}, mongoDBHandler),
	}
}
