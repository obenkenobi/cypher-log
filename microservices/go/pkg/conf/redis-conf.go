package conf

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/environment"
	"github.com/obenkenobi/cypher-log/services/go/pkg/utils"
	"strconv"
)

type RedisConf interface {
	GetAddress() string
	GetPassword() string
	GetDB() int
}

type RedisConfImpl struct {
	addr     string
	password string
	db       int
}

func (r RedisConfImpl) GetAddress() string { return r.addr }

func (r RedisConfImpl) GetPassword() string { return r.password }

func (r RedisConfImpl) GetDB() int { return r.db }

func NewRedisConfImpl(envVarKeyRedisAddr, envVarKeyRedisPassword, envVarKeyRedisDB string) *RedisConfImpl {
	db := 0
	if dbStr := environment.GetEnvVariable(envVarKeyRedisDB); utils.StringIsNotBlank(dbStr) {
		if dbInt, err := strconv.Atoi(dbStr); err == nil {
			db = dbInt
		}
	}
	return &RedisConfImpl{
		addr:     environment.GetEnvVariable(envVarKeyRedisAddr),
		password: environment.GetEnvVariable(envVarKeyRedisPassword),
		db:       db,
	}
}
