package externalservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/concurrent"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"google.golang.org/grpc"
)

type CoreGrpcConnProvider interface {
	CreateConnectionSingle(ctx context.Context, address string) (*grpc.ClientConn, error)
}

type CoreGrpcConnProviderImpl struct {
	systemAccessTokenClient SysAccessTokenClient
	tlsConf                 conf.TLSConf
}

func (u CoreGrpcConnProviderImpl) CreateConnectionSingle(_ context.Context, address string) (*grpc.ClientConn, error) {
	var dialOptions []grpc.DialOption
	if environment.ActivateGRPCAuth() {
		oathTokenFuture := concurrent.Async(u.systemAccessTokenClient.GetGRPCAccessToken)
		if u.tlsConf.WillLoadCACert() {
			tlsOpt, err := gtools.LoadTLSCredentialsOption(u.tlsConf.CACertPath(), environment.IsDevelopment())
			if err != nil {
				return nil, err
			}
			dialOptions = append(dialOptions, tlsOpt)
		}
		token, err := oathTokenFuture.Await()
		if err != nil {
			return nil, err
		}
		dialOptions = append(dialOptions, gtools.OathAccessOption(token))
	}
	return grpc.Dial(address, dialOptions...)
}

func NewCoreGrpcConnProviderImpl(
	systemAccessTokenClient SysAccessTokenClient,
	tlsConf conf.TLSConf,
) *CoreGrpcConnProviderImpl {
	return &CoreGrpcConnProviderImpl{
		systemAccessTokenClient: systemAccessTokenClient,
		tlsConf:                 tlsConf,
	}
}
