package repositories

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/mgmtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
)

type NoteRepository interface {
	baserepos.CRUDRepository[models.Note, string]
	GetPaginatedByUserId(
		ctx context.Context,
		userId string,
		pageReq pagination.PageRequest,
	) ([]models.Note, error)
	CountByUserId(ctx context.Context, userId string) (int64, error)
	DeleteByUserIdAndGetCount(ctx context.Context, userId string) (int64, error)
}

type NoteRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[models.Note]
}

func (u NoteRepositoryImpl) Create(ctx context.Context, model models.Note) (models.Note, error) {
	err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
	return model, err
}

func (u NoteRepositoryImpl) Update(ctx context.Context, model models.Note) (models.Note, error) {
	err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
	return model, err
}

func (u NoteRepositoryImpl) Delete(ctx context.Context, model models.Note) (models.Note, error) {
	err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
	return model, err
}

func (u NoteRepositoryImpl) FindById(ctx context.Context, id string) (option.Maybe[models.Note], error) {
	return dshandlers.HandleSingleFind(u.MongoDBHandler, func() (models.Note, error) {
		model := models.Note{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.ToChildCtx(ctx), id, &model)
		return model, err
	})
}

func (u NoteRepositoryImpl) GetPaginatedByUserId(
	ctx context.Context,
	userId string,
	pageReq pagination.PageRequest,
) ([]models.Note, error) {
	findOpts := mgmtools.CreatePaginatedFindOpts(pageReq)
	filter := bson.D{{"userId", userId}}
	childCtx := u.MongoDBHandler.ToChildCtx(ctx)
	cursor, err := mgm.Coll(u.ModelColl).Find(childCtx, filter, findOpts)
	return mgmtools.HandleFindManyRes[models.Note](childCtx, cursor, err)
}

func (u NoteRepositoryImpl) CountByUserId(ctx context.Context, userId string) (int64, error) {
	filter := bson.D{{"userId", userId}}
	return mgm.Coll(u.ModelColl).CountDocuments(u.MongoDBHandler.ToChildCtx(ctx), filter)
}

func (u NoteRepositoryImpl) DeleteByUserIdAndGetCount(ctx context.Context, userId string) (int64, error) {
	res, err := mgm.Coll(u.ModelColl).DeleteMany(u.MongoDBHandler.ToChildCtx(ctx), bson.M{"userId": userId})
	if res != nil {
		return res.DeletedCount, err
	}
	return -1, err
}

func NewNoteRepositoryImpl(mongoDBHandler *dshandlers.MongoDBHandler) *NoteRepositoryImpl {
	return &NoteRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[models.Note](models.Note{}, mongoDBHandler),
	}
}
