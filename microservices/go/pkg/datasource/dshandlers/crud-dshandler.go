package dshandlers

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
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
	// ToChildCtx creates a new context to be sent to other database queries
	ToChildCtx(ctx context.Context) context.Context
	// ExecTransaction executes a transaction synchronously from the runner function.
	ExecTransaction(ctx context.Context, runner func(Session, context.Context) error) error
}

// Transactional executes a transaction
func Transactional[T any](
	ctx context.Context,
	d CrudDSHandler,
	supplier func(Session, context.Context) (T, error),
) (T, error) {
	var res T
	var err error = nil
	transactionErr := d.ExecTransaction(ctx, func(session Session, ctx context.Context) error {
		res, err = supplier(session, ctx)
		if err != nil {
			return session.AbortTransaction(ctx)
		}
		return session.CommitTransaction(ctx)
	})
	if err == nil {
		err = transactionErr
	}
	return res, err
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

func (d MongoDBHandler) ToChildCtx(ctx context.Context) context.Context {
	dbCtx, _ := d.GetChildDBCtxWithCancel(ctx)
	return dbCtx
}

func (d MongoDBHandler) ExecTransaction(ctx context.Context, transactionFunc func(Session, context.Context) error) error {
	return mgm.TransactionWithCtx(
		ctx,
		func(session mongo.Session, sc mongo.SessionContext) error {
			return transactionFunc(session, sc)
		},
	)
}

func NewMongoDBHandler(mongoConf conf.MongoConf) *MongoDBHandler {
	err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: mongoConf.GetConnectionTimeout()},
		mongoConf.GetDBName(),
		options.Client().ApplyURI(mongoConf.GetUri()))
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to set mongodb config")
	}
	return &MongoDBHandler{mongoConf: mongoConf}
}
