package dbservices

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/containers/option"
	stx "github.com/obenkenobi/cypher-log/services/go/pkg/reactive/single"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session interface {
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
}

// CrudDBHandler Handles CRUD database related tasks such as setting up database
// connection(s), providing contexts for database operations, managing
// transactions, and handling database errors.
type CrudDBHandler interface {
	// GetChildDBCtxWithCancel creates a new context to be sent to other database queries with a cancel function
	GetChildDBCtxWithCancel(ctx context.Context) (context.Context, context.CancelFunc)
	// GetChildDBCtx creates a new context to be sent to other database queries
	GetChildDBCtx(ctx context.Context) context.Context
	// IsNotFoundError checks if an error is created by an underlying object database
	// mapper is due to a requested entity not being found.
	IsNotFoundError(err error) bool
	// ExecTransaction executes a transaction. Warning: not tested and will
	// eventually be scrapped with an implementation meant to work with Singles and
	// Observables.
	ExecTransaction(runner func(Session, context.Context) error) error
}

// MongoDBHandler is a CrudDBHandler implementation for MongoDB, in particular for
// the kamva/mgm ODM.
type MongoDBHandler struct {
	mongoConf conf.MongoConf
}

func (d MongoDBHandler) IsNotFoundError(err error) bool {
	return err.Error() == "mongo: no documents in result"
}

func (d MongoDBHandler) GetChildDBCtxWithCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d.mongoConf.GetConnectionTimeout())
}

func (d MongoDBHandler) GetChildDBCtx(ctx context.Context) context.Context {
	dbCtx, _ := d.GetChildDBCtxWithCancel(ctx)
	return dbCtx
}

func (d MongoDBHandler) ExecTransaction(transactionFunc func(Session, context.Context) error) error {
	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		return transactionFunc(session, sc)
	})
}

func NewMongoHandler(mongoConf conf.MongoConf) *MongoDBHandler {
	if err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: mongoConf.GetConnectionTimeout()},
		mongoConf.GetDBName(),
		options.Client().ApplyURI(mongoConf.GetUri())); err != nil {
		log.WithError(err).Fatal("Failed to set mongodb config")
	}
	return &MongoDBHandler{mongoConf: mongoConf}
}

// ObserveOptionalSingleQueryAsync
//creates a single out of a supplier function that queries a single value. The
//supplier function is run on a separate goroutine. *Make sure your supplier
//function is not going to be thread safe or not cause race conditions on the
//data accessed.
func ObserveOptionalSingleQueryAsync[TQueryResult any](
	mongoDBHandler *MongoDBHandler,
	supplier func() (TQueryResult, error),
) stx.Single[option.Maybe[TQueryResult]] {
	return stx.FromSupplierAsync(func() (option.Maybe[TQueryResult], error) {
		return runOptionalSingleQuery(mongoDBHandler, supplier)
	})
}

func ObserveOptionalSingleQuery[TQueryResult any](
	mongoDBHandler *MongoDBHandler,
	supplier func() (TQueryResult, error),
) stx.Single[option.Maybe[TQueryResult]] {
	return stx.FromSupplier(func() (option.Maybe[TQueryResult], error) {
		return runOptionalSingleQuery(mongoDBHandler, supplier)
	})
}

func runOptionalSingleQuery[TQueryResult any](
	mongoDBHandler *MongoDBHandler,
	supplier func() (TQueryResult, error),
) (option.Maybe[TQueryResult], error) {
	if result, err := supplier(); err != nil {
		if mongoDBHandler.IsNotFoundError(err) {
			return option.None[TQueryResult](), nil
		}
		return nil, err
	} else {
		return option.Perhaps(result), nil
	}
}
