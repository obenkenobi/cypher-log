package grpcserver

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/taskrunner"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

// GrpcServer represents an interface for a grpc server that can be run.
type GrpcServer interface{ taskrunner.TaskRunner }

type grpcServerImpl struct {
	server     *grpc.Server
	serverConf conf.ServerConf
}

func (g grpcServerImpl) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", g.serverConf.GetGrpcServerPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	if err := g.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewGrpcServer(
	serverConf conf.ServerConf,
	registerServers func(server *grpc.Server),
	serverOptions ...grpc.ServerOption,
) GrpcServer {
	server := grpc.NewServer(serverOptions...)
	registerServers(server)
	grpcServerImpl := &grpcServerImpl{serverConf: serverConf, server: server}
	return grpcServerImpl
}
