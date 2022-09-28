package baserepos

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/datasource/dshandlers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type CRUDRepository[VModel any, VID any] interface {
	// Create saves a new model to a data store. The model is updated with the saved
	// values from the database onto the same model and then is emitted by a Single.
	// The model should be a pointer.
	Create(ctx context.Context, modelRef VModel) single.Single[VModel]

	// Update saves an existing model to a data store. The model is updated with the
	// saved values from the data store onto the same model and then is emitted by a
	// Single.
	Update(ctx context.Context, model VModel) single.Single[VModel]

	// Delete deletes an existing model to a data store. The model is then emitted by
	// a Single.
	Delete(ctx context.Context, model VModel) single.Single[VModel]

	// FindById queries the data store by an entity's id and saves the value to the
	// provided model. The same model is then emitted by a Single. The model should
	// be a pointer.
	FindById(ctx context.Context, id VID) single.Single[option.Maybe[VModel]]
}

// BaseRepositoryMongo is a MongoDB implementation of CRUDRepository
type BaseRepositoryMongo[VModel any] struct {
	ModelColl      *VModel
	MongoDBHandler *dshandlers.MongoDBHandler
}

func NewBaseRepositoryMongo[VModel any](modelColl VModel, dbHandler *dshandlers.MongoDBHandler) *BaseRepositoryMongo[VModel] {
	return &BaseRepositoryMongo[VModel]{ModelColl: &modelColl, MongoDBHandler: dbHandler}
}
