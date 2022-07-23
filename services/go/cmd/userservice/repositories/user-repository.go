package repositories

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/cmd/userservice/models"
	"github.com/obenkenobi/cypher-log/services/go/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	FindById(ctx context.Context, id string, user *models.User) error
	FindByAuthId(ctx context.Context, authId string, user *models.User) error
	FindByUsername(ctx context.Context, username string, user *models.User) error
}

type UserRepositoryImpl struct {
	UserColl *models.User
}

func NewUserMongoRepository(mongoDBHandler database.MongoDBHandler) UserRepository {
	return &UserRepositoryImpl{UserColl: &models.User{}}
}

func (u UserRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	return mgm.Coll(u.UserColl).CreateWithCtx(ctx, user)
}

func (u UserRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	return mgm.Coll(u.UserColl).UpdateWithCtx(ctx, user)
}

func (u UserRepositoryImpl) FindById(ctx context.Context, id string, user *models.User) error {
	return mgm.Coll(u.UserColl).FindByIDWithCtx(ctx, id, user)
}

func (u UserRepositoryImpl) FindByAuthId(ctx context.Context, authId string, user *models.User) error {
	return mgm.Coll(u.UserColl).FirstWithCtx(ctx, bson.M{"authId": authId}, user)
}

func (u UserRepositoryImpl) FindByUsername(ctx context.Context, username string, user *models.User) error {
	return mgm.Coll(u.UserColl).FirstWithCtx(ctx, bson.M{"userName": username}, user)
}
