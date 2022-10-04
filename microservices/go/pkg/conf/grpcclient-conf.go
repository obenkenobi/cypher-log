package conf

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type GrpcClientConf interface {
	UserServiceAddress() string
}

type GrpcClientConfImpl struct {
	userServiceAddr string
}

func (g GrpcClientConfImpl) UserServiceAddress() string { return g.userServiceAddr }

func NewGrpcClientConfImpl() *GrpcClientConfImpl {
	return &GrpcClientConfImpl{
		userServiceAddr: environment.GetEnvVariable(environment.EnvVarKeyGrpcUserServiceAddress),
	}
}
