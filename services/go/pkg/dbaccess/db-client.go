package dbaccess

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient interface {
	GetContext() context.Context
}

type DBMongoClient struct {
}

func (D DBMongoClient) GetContext() context.Context {
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
