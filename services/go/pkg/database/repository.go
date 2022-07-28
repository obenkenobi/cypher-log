package database

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/reactorextensions"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
)

type Repository[VModel MongoModel, VID any] interface {
	Create(ctx context.Context, modelRef VModel) stream.Observable[VModel]
	CreateAsync(ctx context.Context, modelRef VModel) stream.Observable[VModel]
	Update(ctx context.Context, modelRef VModel) stream.Observable[VModel]
	UpdateAsync(ctx context.Context, modelRef VModel) stream.Observable[VModel]
	Delete(ctx context.Context, modelRef VModel) stream.Observable[VModel]
	DeleteAsync(ctx context.Context, modelRef VModel) stream.Observable[VModel]
	FindById(ctx context.Context, modelRef VModel, id string) stream.Observable[option.Maybe[VModel]]
	FindByIdAsync(ctx context.Context, modelRef VModel, id string) stream.Observable[option.Maybe[VModel]]
}

type RepositoryMongoImpl[VModel MongoModel, VID any] struct {
	ModelColumn    VModel
	MongoDBHandler *MongoDBHandler
}

func (r RepositoryMongoImpl[VModel, VID]) Create(ctx context.Context, modelRef VModel) stream.Observable[VModel] {
	return reactorextensions.ObserveSupplier(func() (VModel, error) {
		return r.runCreate(ctx, modelRef)
	})
}

func (r RepositoryMongoImpl[VModel, VID]) CreateAsync(ctx context.Context, modelRef VModel) stream.Observable[VModel] {
	return reactorextensions.ObserveSupplierAsync(func() (VModel, error) {
		return r.runCreate(ctx, modelRef)
	})
}

func (r RepositoryMongoImpl[VModel, VID]) runCreate(ctx context.Context, modelRef VModel) (VModel, error) {
	err := mgm.Coll(r.ModelColumn).CreateWithCtx(ctx, modelRef)
	return modelRef, err
}

func (r RepositoryMongoImpl[VModel, VID]) Update(ctx context.Context, modelRef VModel) stream.Observable[VModel] {
	return reactorextensions.ObserveSupplier(func() (VModel, error) {
		return r.runUpdate(ctx, modelRef)
	})
}

func (r RepositoryMongoImpl[VModel, VID]) UpdateAsync(ctx context.Context, modelRef VModel) stream.Observable[VModel] {
	return reactorextensions.ObserveSupplierAsync(func() (VModel, error) {
		return r.runUpdate(ctx, modelRef)
	})
}

func (r RepositoryMongoImpl[VModel, VID]) runUpdate(ctx context.Context, modelRef VModel) (VModel, error) {
	err := mgm.Coll(r.ModelColumn).UpdateWithCtx(ctx, modelRef)
	return modelRef, err
}

func (r RepositoryMongoImpl[VModel, VID]) Delete(ctx context.Context, modelRef VModel) stream.Observable[VModel] {
	return reactorextensions.ObserveSupplier(func() (VModel, error) {
		return r.runDelete(ctx, modelRef)
	})
}

func (r RepositoryMongoImpl[VModel, VID]) DeleteAsync(ctx context.Context, modelRef VModel) stream.Observable[VModel] {
	return reactorextensions.ObserveSupplierAsync(func() (VModel, error) {
		return r.runDelete(ctx, modelRef)
	})
}

func (r RepositoryMongoImpl[VModel, VID]) runDelete(ctx context.Context, modelRef VModel) (VModel, error) {
	err := mgm.Coll(r.ModelColumn).DeleteWithCtx(ctx, modelRef)
	return modelRef, err
}

func (r RepositoryMongoImpl[VModel, VID]) FindById(
	ctx context.Context,
	modelRef VModel,
	id string,
) stream.Observable[option.Maybe[VModel]] {
	return ObserveOptionalSingleQueryAsync(r.MongoDBHandler, func() (VModel, error) {
		return r.runFindByIdAsync(ctx, modelRef, id)
	})
}

func (r RepositoryMongoImpl[VModel, VID]) FindByIdAsync(
	ctx context.Context,
	modelRef VModel,
	id string,
) stream.Observable[option.Maybe[VModel]] {
	return ObserveOptionalSingleQueryAsync(r.MongoDBHandler, func() (VModel, error) {
		return r.runFindByIdAsync(ctx, modelRef, id)
	})
}

func (r RepositoryMongoImpl[VModel, VID]) runFindByIdAsync(
	ctx context.Context,
	modelRef VModel,
	id string,
) (VModel, error) {
	err := mgm.Coll(r.ModelColumn).FindByIDWithCtx(ctx, id, modelRef)
	return modelRef, err
}

func NewRepositoryMongoImpl[VModel MongoModel, VID any](
	modelColumn VModel,
	mongoDBHandler *MongoDBHandler,
) *RepositoryMongoImpl[VModel, VID] {
	return &RepositoryMongoImpl[VModel, VID]{ModelColumn: modelColumn, MongoDBHandler: mongoDBHandler}
}
