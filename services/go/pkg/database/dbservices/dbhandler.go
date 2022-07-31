package dbservices

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	stx "github.com/obenkenobi/cypher-log/services/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/services/go/pkg/wrappers/option"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session interface {
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
}

// DBHandler Handles database related tasks such as setting up database connection(s),
// providing contexts for database operations, managing transactions, and handling database apperrors.
type DBHandler interface {
	GetCtx() context.Context
	IsNotFoundError(err error) bool
	ExecTransaction(runner func(Session, context.Context) error) error
}

type MongoDBHandler struct {
}

func (d MongoDBHandler) IsNotFoundError(err error) bool {
	return err.Error() == "mongo: no documents in result"
}

func (d MongoDBHandler) GetCtx() context.Context {
	return mgm.Ctx()
}

func (d MongoDBHandler) ExecTransaction(transactionFunc func(Session, context.Context) error) error {
	return mgm.Transaction(func(session mongo.Session, sc mongo.SessionContext) error {
		return transactionFunc(session, sc)
	})
}

func BuildMongoHandler(mongoConf conf.MongoConf) *MongoDBHandler {
	if err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: mongoConf.GetConnectionTimeout()},
		mongoConf.GetDBName(),
		options.Client().ApplyURI(mongoConf.GetUri())); err != nil {
		log.WithError(err).Fatal("Failed to set mongodb config")
	}
	return &MongoDBHandler{}
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
