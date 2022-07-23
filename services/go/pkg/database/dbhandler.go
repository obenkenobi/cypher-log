package database

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session interface {
	AbortTransaction(context.Context) error
	CommitTransaction(context.Context) error
}

// DBHandler Handles database related tasks such as setting up database connection(s),
// providing contexts for database operations, managing transactions, and handling database errors.
type DBHandler interface {
	GetCtx() context.Context
	IsNotFoundError(err error) bool
	NotFoundOrElseInternalErrResponse(err error) *errordtos.ErrorResponseDto
	ExecTransaction(runner func(Session, context.Context) error) error
}

type MongoDBHandler struct {
}

func (d MongoDBHandler) NotFoundOrElseInternalErrResponse(err error) *errordtos.ErrorResponseDto {
	if d.IsNotFoundError(err) {
		return apperrors.CreateErrorResponseFromErrorCodes(apperrors.ErrCodeReqItemsNotFound)
	}
	return apperrors.CreateInternalErrResponse(err)
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

func BuildMongoHandler(mongoConf conf.MongoConf) MongoDBHandler {
	if err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: mongoConf.GetConnectionTimeout()},
		mongoConf.GetDBName(),
		options.Client().ApplyURI(mongoConf.GetUri())); err != nil {
		log.WithError(err).Fatal("Failed to set mongodb config")
	}
	return MongoDBHandler{}
}
