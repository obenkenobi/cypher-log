package conf

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"

type GrpcClientConf interface {
	UserServiceAddress() string
	KeyServiceAddress() string
}

type GrpcClientConfImpl struct {
	userServiceAddr string
	keyServiceAddr  string
}

func (g GrpcClientConfImpl) UserServiceAddress() string {
	return g.userServiceAddr
}

func (g GrpcClientConfImpl) KeyServiceAddress() string {
	return g.keyServiceAddr
}

func NewGrpcClientConfImpl() *GrpcClientConfImpl {
	return &GrpcClientConfImpl{
		userServiceAddr: environment.GetEnvVariable(environment.EnvVarKeyGrpcUserServiceAddress),
		keyServiceAddr:  environment.GetEnvVariable(environment.EnvVarKeyGrpcKeyServiceAddress),
	}
}
