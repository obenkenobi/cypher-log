package conf

import (
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

func NewMongoConf(envVarReader EnvVarReader, envVarKeyMongoUri string, envVarMongoDBName string,
	envVarMongoConnTimeoutMS string) MongoConf {
	connTimeout := 12 * time.Second
	if connTimeoutStr := envVarReader.GetEnvVariable(envVarMongoConnTimeoutMS); connTimeoutStr != "" {
		if connectionTimeoutInt, err := strconv.ParseInt(connTimeoutStr, 10, 64); err == nil {
			connTimeout = time.Duration(connectionTimeoutInt)
		}
	}

	return &MongoConfImpl{
		mongoUri:          envVarReader.GetEnvVariable(envVarKeyMongoUri),
		mongoDBName:       envVarReader.GetEnvVariable(envVarMongoDBName),
		connectionTimeout: connTimeout,
	}
}
