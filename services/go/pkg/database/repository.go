package database

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/reactorextensions"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
)

type Repository[VModel MongoModel, VID any] interface {
	Create(ctx context.Context, model VModel) stream.Observable[VModel]
	Update(ctx context.Context, model VModel) stream.Observable[VModel]
	Delete(ctx context.Context, model VModel) stream.Observable[VModel]
	FindById(ctx context.Context, model VModel, id string) stream.Observable[option.Maybe[VModel]]
}

type RepositoryMongoImpl[VModel MongoModel, VID any] struct {
	ModelColumn    VModel
	MongoDBHandler *MongoDBHandler
}

func (r RepositoryMongoImpl[VModel, VID]) Create(ctx context.Context, model VModel) stream.Observable[VModel] {
	return reactorextensions.ObserveProducer(func() (VModel, error) {
		err := mgm.Coll(r.ModelColumn).CreateWithCtx(ctx, model)
		return model, err
	})
}

func (r RepositoryMongoImpl[VModel, VID]) Update(ctx context.Context, model VModel) stream.Observable[VModel] {
	return reactorextensions.ObserveProducer(func() (VModel, error) {
		err := mgm.Coll(r.ModelColumn).UpdateWithCtx(ctx, model)
		return model, err
	})
}

func (r RepositoryMongoImpl[VModel, VID]) Delete(ctx context.Context, model VModel) stream.Observable[VModel] {
	return reactorextensions.ObserveProducer(func() (VModel, error) {
		err := mgm.Coll(r.ModelColumn).DeleteWithCtx(ctx, model)
		return model, err
	})
}

func (r RepositoryMongoImpl[VModel, VID]) FindById(
	ctx context.Context,
	model VModel,
	id string,
) stream.Observable[option.Maybe[VModel]] {
	return ObserveOptionalSingleQuery(r.MongoDBHandler, func() (VModel, error) {
		err := mgm.Coll(r.ModelColumn).FindByIDWithCtx(ctx, id, model)
		return model, err
	})
}

func NewRepositoryMongoImpl[VModel MongoModel, VID any](
	modelColumn VModel,
	mongoDBHandler *MongoDBHandler,
) *RepositoryMongoImpl[VModel, VID] {
	return &RepositoryMongoImpl[VModel, VID]{ModelColumn: modelColumn, MongoDBHandler: mongoDBHandler}
}
