package conf

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
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

func NewMongoConfImpl() *MongoConfImpl {
	connTimeout := environment.GetEnvVarAsTimeDurationOrDefault(environment.EnvVarMongoConnTimeoutMS, 12*time.Second)
	return &MongoConfImpl{
		mongoUri:          environment.GetEnvVar(environment.EnvVarKeyMongoUri),
		mongoDBName:       environment.GetEnvVar(environment.EnvVarMongoDBName),
		connectionTimeout: connTimeout,
	}
}
