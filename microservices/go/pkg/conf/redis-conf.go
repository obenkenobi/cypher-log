package conf

import (
	environment2 "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
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

func NewRedisConfImpl() *RedisConfImpl {
	return &RedisConfImpl{
		addr:     environment2.GetEnvVariable(environment2.EnvVarKeyRedisAddr),
		password: environment2.GetEnvVariable(environment2.EnvVarKeyRedisPassword),
		db:       environment2.GetEnvVarAsIntOrDefault(environment2.EnvVarKeyRedisDB, 0),
	}
}
