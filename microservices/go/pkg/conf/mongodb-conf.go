package conf

import (
	environment2 "github.com/obenkenobi/cypher-log/services/go/pkg/environment"
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

func NewMongoConf() MongoConf {
	connTimeout := environment2.GetEnvVarAsTimeDurationOrDefault(environment2.EnvVarMongoConnTimeoutMS, 12*time.Second)
	return &MongoConfImpl{
		mongoUri:          environment2.GetEnvVariable(environment2.EnvVarKeyMongoUri),
		mongoDBName:       environment2.GetEnvVariable(environment2.EnvVarMongoDBName),
		connectionTimeout: connTimeout,
	}
}
