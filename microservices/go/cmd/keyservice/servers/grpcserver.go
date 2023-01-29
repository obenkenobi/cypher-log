package servers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/commonservers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userkeypb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
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
	userKeyServiceServer userkeypb.UserKeyServiceServer,
) *GrpcServerImpl {
	if !environment.ActivateGrpcServer() {
		// Server is deactivated, ran via the lifecycle package,
		// and is a root-child dependency so a nil is returned
		return nil
	}
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
			userkeypb.RegisterUserKeyServiceServer(s, userKeyServiceServer)
		},
		grpcOpts...,
	)
	g := &GrpcServerImpl{CoreGrpcServer: coreServer}
	lifecycle.RegisterTaskRunner(g)
	return g
}
