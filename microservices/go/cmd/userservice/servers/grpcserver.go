package servers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices/grpcserveropts"
	"google.golang.org/grpc"
)

type GrpcServer interface {
	commonservers.CoreGrpcServer
}

type GrpcServerImpl struct {
	commonservers.CoreGrpcServer
	authInterceptorCreator   grpcserveropts.AuthInterceptorCreator
	credentialsOptionCreator grpcserveropts.CredentialsOptionCreator
}

func NewGrpcServerImpl(
	serverConf conf.ServerConf,
	authInterceptorCreator grpcserveropts.AuthInterceptorCreator,
	credentialsOptionCreator grpcserveropts.CredentialsOptionCreator,
	userServiceServer userpb.UserServiceServer,
) *GrpcServerImpl {
	var grpcOpts []grpc.ServerOption
	if environment.ActivateGRPCAuth() {
		grpcOpts = append(
			grpcOpts,
			authInterceptorCreator.CreateUnaryInterceptor(),
			credentialsOptionCreator.CreateCredentialsOption(),
		)
	}
	coreServer := commonservers.NewCoreGrpcServer(
		serverConf,
		func(s *grpc.Server) {
			userpb.RegisterUserServiceServer(s, userServiceServer)
		},
		grpcOpts...,
	)
	return &GrpcServerImpl{CoreGrpcServer: coreServer}
}
