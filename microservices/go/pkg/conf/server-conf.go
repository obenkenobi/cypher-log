package conf

import (
	environment2 "github.com/obenkenobi/cypher-log/services/go/pkg/environment"
)

type ServerConf interface {
	GetAppServerPort() string
	GetGrpcServerPort() string
}

type serverConfImpl struct {
	appServerPort  string
	grpcServerPort string
}

func (s serverConfImpl) GetGrpcServerPort() string { return s.grpcServerPort }

func (s serverConfImpl) GetAppServerPort() string { return s.appServerPort }

func NewServerConf() ServerConf {
	return &serverConfImpl{
		appServerPort:  environment2.GetEnvVariableOrDefault(environment2.EnvVarKeyAppServerPort, "8080"),
		grpcServerPort: environment2.GetEnvVariableOrDefault(environment2.EnvVarKeyGrpcServerPort, "5000"),
	}
}
