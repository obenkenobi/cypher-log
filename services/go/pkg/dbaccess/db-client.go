package dbaccess

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"
	"github.com/obenkenobi/cypher-log/services/go/pkg/errormgmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient interface {
	GetCtx() context.Context
	IsNotFoundError(err error) bool
	NotFoundOrElseInternalErrResponse(err error) *errordtos.ErrorResponseDto
}

type DBMongoClient struct {
}

func (d DBMongoClient) NotFoundOrElseInternalErrResponse(err error) *errordtos.ErrorResponseDto {
	if d.IsNotFoundError(err) {
		return errormgmt.CreateErrorResponseFromErrorCodes(errormgmt.ErrCodeReqItemsNotFound)
	}
	return errormgmt.CreateInternalErrResponseWithErrLog(err)
}

func (D DBMongoClient) IsNotFoundError(err error) bool {
	return err.Error() == "mongo: no documents in result"
}

func (D DBMongoClient) GetCtx() context.Context {
	return mgm.Ctx()
}

func BuildMongoClient(mongoConf conf.MongoConf) DBMongoClient {
	if err := mgm.SetDefaultConfig(
		&mgm.Config{CtxTimeout: mongoConf.GetConnectionTimeout()},
		mongoConf.GetDBName(),
		options.Client().ApplyURI(mongoConf.GetUri())); err != nil {
		log.WithError(err).Fatal("Failed to set mongodb config")
	}

	return DBMongoClient{}
}
