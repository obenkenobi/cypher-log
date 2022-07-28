package conf

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"
	"strconv"
	"time"
)

type MongoConf interface {
	GetUri() string
	GetDBName() string
	GetConnectionTimeout() time.Duration
}

type MongoConfImpl struct {
	mongoUri          string
	mongoDBName       string
	connectionTimeout time.Duration
}

func (m MongoConfImpl) GetUri() string {
	return m.mongoUri
}

func (m MongoConfImpl) GetDBName() string {
	return m.mongoDBName
}

func (m MongoConfImpl) GetConnectionTimeout() time.Duration {
	return m.connectionTimeout
}

func NewMongoConf(envVarKeyMongoUri string, envVarKeyMongoDBName string,
	envVarKeyMongoConnTimeoutMS string) MongoConf {
	connTimeout := 12 * time.Second
	if connTimeoutStr := environment.GetEnvVariable(envVarKeyMongoConnTimeoutMS); connTimeoutStr != "" {
		if connectionTimeoutInt, err := strconv.ParseInt(connTimeoutStr, 10, 64); err == nil {
			connTimeout = time.Duration(connectionTimeoutInt)
		}
	}

	return &MongoConfImpl{
		mongoUri:          environment.GetEnvVariable(envVarKeyMongoUri),
		mongoDBName:       environment.GetEnvVariable(envVarKeyMongoDBName),
		connectionTimeout: connTimeout,
	}
}
