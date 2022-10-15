package repositories

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/baserepos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type NoteRepository interface {
	baserepos.CRUDRepository[models.Note, string]
}

type NoteRepositoryImpl struct {
	baserepos.BaseRepositoryMongo[models.Note]
}

func (u NoteRepositoryImpl) Create(
	ctx context.Context,
	model models.Note,
) single.Single[models.Note] {
	return single.FromSupplierCached(func() (models.Note, error) {
		err := mgm.Coll(u.ModelColl).CreateWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &model)
		return model, err
	})
}

func (u NoteRepositoryImpl) Update(
	ctx context.Context,
	model models.Note,
) single.Single[models.Note] {
	return single.FromSupplierCached(func() (models.Note, error) {
		err := mgm.Coll(u.ModelColl).UpdateWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &model)
		return model, err
	})
}

func (u NoteRepositoryImpl) Delete(
	ctx context.Context,
	model models.Note,
) single.Single[models.Note] {
	return single.FromSupplierCached(func() (models.Note, error) {
		err := mgm.Coll(u.ModelColl).DeleteWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), &model)
		return model, err
	})
}

func (u NoteRepositoryImpl) FindById(
	ctx context.Context,
	id string,
) single.Single[option.Maybe[models.Note]] {
	return dshandlers.OptionalSingleQuerySrc(u.MongoDBHandler, func() (models.Note, error) {
		model := models.Note{}
		err := mgm.Coll(u.ModelColl).FindByIDWithCtx(u.MongoDBHandler.GetChildDBCtx(ctx), id, &model)
		return model, err
	})
}

func NewNoteRepositoryImpl(mongoDBHandler *dshandlers.MongoDBHandler) *NoteRepositoryImpl {
	return &NoteRepositoryImpl{
		BaseRepositoryMongo: *baserepos.NewBaseRepositoryMongo[models.Note](
			models.Note{},
			mongoDBHandler,
		),
	}
}
