package repositories

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/mgmtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/pagination"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"go.mongodb.org/mongo-driver/bson"
)

type NoteRepository interface {
	baserepos.CRUDRepository[models.Note, string]
	FindManyByUserId(ctx context.Context, userId string, pageReq pagination.PageRequest) stream.Observable[models.Note]
}

type NoteRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[models.Note]
}

func (u NoteRepositoryImpl) Create(ctx context.Context, model models.Note) single.Single[models.Note] {
	return single.FromSupplierCached(func() (models.Note, error) {
		err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
		return model, err
	})
}

func (u NoteRepositoryImpl) Update(ctx context.Context, model models.Note) single.Single[models.Note] {
	return single.FromSupplierCached(func() (models.Note, error) {
		err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
		return model, err
	})
}

func (u NoteRepositoryImpl) Delete(ctx context.Context, model models.Note) single.Single[models.Note] {
	return single.FromSupplierCached(func() (models.Note, error) {
		err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.ToChildCtx(ctx), &model)
		return model, err
	})
}

func (u NoteRepositoryImpl) FindById(ctx context.Context, id string) single.Single[option.Maybe[models.Note]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.Note, error) {
		model := models.Note{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.ToChildCtx(ctx), id, &model)
		return model, err
	})
}

func (u NoteRepositoryImpl) FindManyByUserId(
	ctx context.Context,
	userId string,
	pageReq pagination.PageRequest,
) stream.Observable[models.Note] {
	return stream.FlatMap(stream.Just([]models.Note{}), func(results []models.Note) stream.Observable[models.Note] {
		findOpts := mgmtools.CreatePaginatedFindOpts(pageReq)
		filter := bson.D{{"userId", userId}}
		cursor, err := mgm.Coll(u.ModelColl).Find(u.MongoDBHandler.ToChildCtx(ctx), filter, findOpts)
		return mgmtools.HandleFindManyRes(ctx, results, cursor, err)
	})
}

func NewNoteRepositoryImpl(mongoDBHandler *dshandlers.MongoDBHandler) *NoteRepositoryImpl {
	return &NoteRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[models.Note](models.Note{}, mongoDBHandler),
	}
}
