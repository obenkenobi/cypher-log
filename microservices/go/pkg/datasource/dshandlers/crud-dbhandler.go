package dshandlers

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session interface {
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
}

// CrudDSHandler Handles CRUD data source related tasks such as setting up database
// connection(s), providing contexts for database operations, managing
// transactions, and handling database errors.
type CrudDSHandler interface {
	DataSourceHandler
	// GetChildDBCtxWithCancel creates a new context to be sent to other database queries with a cancel function
	GetChildDBCtxWithCancel(ctx context.Context) (context.Context, context.CancelFunc)
	// GetChildDBCtx creates a new context to be sent to other database queries
	GetChildDBCtx(ctx context.Context) context.Context
	// ExecTransaction executes a transaction synchronously from the runner function.
	ExecTransaction(runner func(Session, context.Context) error) error
}

// MongoDBHandler is a CrudDSHandler implementation for MongoDB
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

// TransactionalSingle creates a Single that executes a transaction when
// evaluated from a Single created from the supplier. The supplier and the
// evaluation of the single runs within the scope of the transaction.
func TransactionalSingle[T any](
	d CrudDSHandler,
	supplier func(Session, context.Context) single.Single[T],
) single.Single[T] {
	return single.FromSupplier(func() (T, error) {
		var res T
		var err error = nil
		transactionErr := d.ExecTransaction(func(session Session, ctx context.Context) error {
			res, err = single.RetrieveValue(ctx, supplier(session, ctx))
			if err != nil {
				return session.AbortTransaction(ctx)
			}
			return session.CommitTransaction(ctx)
		})
		if err == nil {
			err = transactionErr
		}
		return res, err
	})
}

// TransactionalObservable creates a deferred Observable that waits for
// a transaction to be completed. The supplier function runs within the
// transaction scope. The returned observable from the supplier is evaluated
// asynchronously and eagerly within the transaction scope.
func TransactionalObservable[T any](
	d CrudDSHandler,
	supplier func(Session, context.Context) stream.Observable[T],
) stream.Observable[T] {
	src, start := stream.Deferred[T]()
	go func() {
		var res []T
		var err error
		transactionErr := d.ExecTransaction(func(session Session, ctx context.Context) error {
			res, err = stream.ToSlice(ctx, supplier(session, ctx))
			if err != nil {
				return session.AbortTransaction(ctx)
			}
			return session.CommitTransaction(ctx)
		})
		if err == nil {
			err = transactionErr
		}
		if err != nil {
			start(stream.Error[T](err))
		}
		start(stream.FromSlice(res))
	}()
	return src
}

func NewMongoHandler(mongoConf conf.MongoConf) *MongoDBHandler {
	err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: mongoConf.GetConnectionTimeout()},
		mongoConf.GetDBName(),
		options.Client().ApplyURI(mongoConf.GetUri()))
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to set mongodb config")
	}
	return &MongoDBHandler{mongoConf: mongoConf}
}
