package database

import (
	"context"
	"github.com/joamaki/goreactive/stream"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/reactorextensions"
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

func ObserveOptionalSingleQueryAsync[TModel any](
	mongoDBHandler *MongoDBHandler,
	producer func() (TModel, error),
) stream.Observable[option.Maybe[TModel]] {
	return reactorextensions.ObserveSupplierAsync(func() (option.Maybe[TModel], error) {
		return runOptionalSingleQuery(mongoDBHandler, producer)
	})
}

func ObserveOptionalSingleQuery[TModel any](
	mongoDBHandler *MongoDBHandler,
	producer func() (TModel, error),
) stream.Observable[option.Maybe[TModel]] {
	return reactorextensions.ObserveSupplier(func() (option.Maybe[TModel], error) {
		return runOptionalSingleQuery(mongoDBHandler, producer)
	})
}

func runOptionalSingleQuery[TModel any](
	mongoDBHandler *MongoDBHandler,
	producer func() (TModel, error),
) (option.Maybe[TModel], error) {
	if result, err := producer(); err != nil {
		if mongoDBHandler.IsNotFoundError(err) {
			return option.None[TModel](), nil
		}
		return nil, err
	} else {
		return option.Perhaps(result), nil
	}
}
