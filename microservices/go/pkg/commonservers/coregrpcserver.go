package commonservers

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/lifecycle"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"google.golang.org/grpc"
	"net"
)

// CoreGrpcServer represents an interface for a grpc server that can be run.
type CoreGrpcServer interface{ lifecycle.TaskRunner }

type coreGrpcServerImpl struct {
	server     *grpc.Server
	serverConf conf.ServerConf
}

func (g coreGrpcServerImpl) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", g.serverConf.GetGrpcServerPort()))
	if err != nil {
		logger.Log.Fatalf("failed to listen: %v", err)
	}
	logger.Log.Infof("server listening at %v", lis.Addr())
	if err := g.server.Serve(lis); err != nil {
		logger.Log.Fatalf("failed to serve: %v", err)
	}
}

func NewCoreGrpcServer(
	serverConf conf.ServerConf,
	registerServers func(server *grpc.Server),
	serverOptions ...grpc.ServerOption,
) CoreGrpcServer {
	server := grpc.NewServer(serverOptions...)
	registerServers(server)
	grpcServerImpl := &coreGrpcServerImpl{serverConf: serverConf, server: server}
	return grpcServerImpl
}
