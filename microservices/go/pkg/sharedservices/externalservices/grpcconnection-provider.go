package externalservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/concurrent"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"google.golang.org/grpc"
)

type CoreGrpcConnProvider interface {
	CreateConnectionSingle(ctx context.Context, address string) single.Single[*grpc.ClientConn]
}

type CoreGrpcConnProviderImpl struct {
	systemAccessTokenClient SysAccessTokenClient
	tlsConf                 conf.TLSConf
}

func (u CoreGrpcConnProviderImpl) CreateConnectionSingle(_ context.Context, address string) single.Single[*grpc.ClientConn] {
	var dialOptions []grpc.DialOption
	if environment.ActivateGRPCAuth() {
		oathTokenFuture := concurrent.Async(u.systemAccessTokenClient.GetGRPCAccessToken)
		if u.tlsConf.WillLoadCACert() {
			tlsOpt, err := gtools.LoadTLSCredentialsOption(u.tlsConf.CACertPath(), environment.IsDevelopment())
			if err != nil {
				return single.Error[*grpc.ClientConn](err)
			}
			dialOptions = append(dialOptions, tlsOpt)
		}
		token, err := oathTokenFuture.Await()
		if err != nil {
			return single.Error[*grpc.ClientConn](err)
		}
		dialOptions = append(dialOptions, gtools.OathAccessOption(token))
	}
	return single.FromSupplierCached(func() (*grpc.ClientConn, error) {
		return grpc.Dial(address, dialOptions...)
	})
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
