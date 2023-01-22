package conf

import (
	environment2 "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
)

type ServerConf interface {
	GetAppServerPort() string
	GetGrpcServerPort() string
}

type ServerConfImpl struct {
	appServerPort  string
	grpcServerPort string
}

func (s ServerConfImpl) GetGrpcServerPort() string { return s.grpcServerPort }

func (s ServerConfImpl) GetAppServerPort() string { return s.appServerPort }

func NewServerConfImpl() *ServerConfImpl {
	return &ServerConfImpl{
		appServerPort:  environment2.GetEnvVarOrDefault(environment2.EnvVarKeyAppServerPort, "8080"),
		grpcServerPort: environment2.GetEnvVarOrDefault(environment2.EnvVarKeyGrpcServerPort, "50051"),
	}
}
