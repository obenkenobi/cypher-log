package dbservices

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/database"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type CRUDRepository[VModel database.MongoModel, VID any] interface {
	// Create saves a new model to a data store. The model is updated with the saved
	// values from the database onto the same model and then is emitted by a Single.
	// The model should be a pointer.
	Create(ctx context.Context, modelRef VModel) single.Single[VModel]

	// CreateAsync saves a new model to a data store. The model is updated with the
	// saved values from the data store onto the same model and then is emitted by a
	// Single. The operation is asynchronous (i.e. runs on another goroutine). The
	// model should be a pointer.
	CreateAsync(ctx context.Context, modelRef VModel) single.Single[VModel]

	// Update saves an existing model to a data store. The model is updated with the
	// saved values from the data store onto the same model and then is emitted by a
	// Single. The model should be a pointer.
	Update(ctx context.Context, modelRef VModel) single.Single[VModel]

	// UpdateAsync saves an existing model to a data store. The model is updated with
	// the saved values from the data store onto the same model and then is emitted
	// by a Single. The operation is asynchronous (i.e. runs on another goroutine).
	// The model should be a pointer.
	UpdateAsync(ctx context.Context, modelRef VModel) single.Single[VModel]

	// Delete deletes an existing model to a data store. The model is then emitted by
	// a Single. The model should be a pointer.
	Delete(ctx context.Context, modelRef VModel) single.Single[VModel]

	// DeleteAsync deletes an existing model to a data store. The model is then
	// emitted by a Single. The operation is asynchronous (i.e. runs on another
	// goroutine). The model should be a pointer.
	DeleteAsync(ctx context.Context, modelRef VModel) single.Single[VModel]

	// FindById queries the data store by an entity's id and saves the value to the
	// provided model. The same model is then emitted by a Single. The model should
	// be a pointer.
	FindById(ctx context.Context, modelRef VModel, id string) single.Single[option.Maybe[VModel]]

	// FindByIdAsync queries the data store by an entity's id and saves the value to
	// the provided model. The operation is asynchronous (i.e. runs on another
	// goroutine). The same model is then emitted by a Single. The model should be a
	// pointer.
	FindByIdAsync(ctx context.Context, modelRef VModel, id string) single.Single[option.Maybe[VModel]]
}

// CRUDRepositoryMongo is a MongoDB implementation of CRUDRepository
type CRUDRepositoryMongo[VModel database.MongoModel, VID any] struct {
	ModelColumn    VModel
	MongoDBHandler *MongoDBHandler
}

func (r CRUDRepositoryMongo[VModel, VID]) Create(ctx context.Context, modelRef VModel) single.Single[VModel] {
	return single.FromSupplier(func() (VModel, error) { return r.runCreate(ctx, modelRef) })
}

func (r CRUDRepositoryMongo[VModel, VID]) CreateAsync(ctx context.Context, modelRef VModel) single.Single[VModel] {
	return single.FromSupplierAsync(func() (VModel, error) { return r.runCreate(ctx, modelRef) })
}

func (r CRUDRepositoryMongo[VModel, VID]) runCreate(ctx context.Context, modelRef VModel) (VModel, error) {
	err := mgm.Coll(r.ModelColumn).CreateWithCtx(r.MongoDBHandler.GetChildDBCtx(ctx), modelRef)
	return modelRef, err
}

func (r CRUDRepositoryMongo[VModel, VID]) Update(ctx context.Context, modelRef VModel) single.Single[VModel] {
	return single.FromSupplier(func() (VModel, error) { return r.runUpdate(ctx, modelRef) })
}

func (r CRUDRepositoryMongo[VModel, VID]) UpdateAsync(ctx context.Context, modelRef VModel) single.Single[VModel] {
	return single.FromSupplierAsync(func() (VModel, error) { return r.runUpdate(ctx, modelRef) })
}

func (r CRUDRepositoryMongo[VModel, VID]) runUpdate(ctx context.Context, modelRef VModel) (VModel, error) {
	err := mgm.Coll(r.ModelColumn).UpdateWithCtx(r.MongoDBHandler.GetChildDBCtx(ctx), modelRef)
	return modelRef, err
}

func (r CRUDRepositoryMongo[VModel, VID]) Delete(ctx context.Context, modelRef VModel) single.Single[VModel] {
	return single.FromSupplier(func() (VModel, error) { return r.runDelete(ctx, modelRef) })
}

func (r CRUDRepositoryMongo[VModel, VID]) DeleteAsync(ctx context.Context, modelRef VModel) single.Single[VModel] {
	return single.FromSupplierAsync(func() (VModel, error) { return r.runDelete(ctx, modelRef) })
}

func (r CRUDRepositoryMongo[VModel, VID]) runDelete(ctx context.Context, modelRef VModel) (VModel, error) {
	err := mgm.Coll(r.ModelColumn).DeleteWithCtx(r.MongoDBHandler.GetChildDBCtx(ctx), modelRef)
	return modelRef, err
}

func (r CRUDRepositoryMongo[VModel, VID]) FindById(
	ctx context.Context,
	modelRef VModel,
	id string,
) single.Single[option.Maybe[VModel]] {
	return ObserveOptionalSingleQueryAsync(r.MongoDBHandler, func() (VModel, error) {
		return r.runFindByIdAsync(ctx, modelRef, id)
	})
}

func (r CRUDRepositoryMongo[VModel, VID]) FindByIdAsync(
	ctx context.Context,
	modelRef VModel,
	id string,
) single.Single[option.Maybe[VModel]] {
	return ObserveOptionalSingleQueryAsync(r.MongoDBHandler, func() (VModel, error) {
		return r.runFindByIdAsync(ctx, modelRef, id)
	})
}

func (r CRUDRepositoryMongo[VModel, VID]) runFindByIdAsync(
	ctx context.Context,
	modelRef VModel,
	id string,
) (VModel, error) {
	err := mgm.Coll(r.ModelColumn).FindByIDWithCtx(r.MongoDBHandler.GetChildDBCtx(ctx), id, modelRef)
	return modelRef, err
}

func NewRepositoryMongoImpl[VModel database.MongoModel, VID any](
	modelColumn VModel,
	mongoDBHandler *MongoDBHandler,
) *CRUDRepositoryMongo[VModel, VID] {
	return &CRUDRepositoryMongo[VModel, VID]{ModelColumn: modelColumn, MongoDBHandler: mongoDBHandler}
}
