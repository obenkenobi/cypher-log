package webservices

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type GrpcServer interface{ Server }

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
