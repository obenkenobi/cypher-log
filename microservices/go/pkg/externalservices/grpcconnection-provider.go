package externalservices

import (
	"context"
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

func (u CoreGrpcConnProviderImpl) CreateConnectionSingle(ctx context.Context, address string) single.Single[*grpc.ClientConn] {
	var dialOptSources []single.Single[grpc.DialOption]
	if environment.ActivateGRPCAuth() {
		oathTokenSrc := single.FromSupplierAsync(u.systemAccessTokenClient.GetGRPCAccessToken)
		tlsOptSrc := single.FromSupplierAsync(func() (grpc.DialOption, error) {
			return gtools.LoadTLSCredentialsOption(u.tlsConf.CACertPath(), environment.IsDevelopment())
		})
		oathOptSrc := single.Map(oathTokenSrc, gtools.OathAccessOption)
		dialOptSources = append(dialOptSources, oathOptSrc, tlsOptSrc)
	}
	optsSrc := gtools.CreateSingleWithDialOptions(dialOptSources)
	return single.MapWithError(optsSrc, func(opts []grpc.DialOption) (*grpc.ClientConn, error) {
		return gtools.CreateConnectionWithOptions(address, opts...)
	})
}

func NewCoreGrpcConnProvider(
	systemAccessTokenClient SysAccessTokenClient,
	tlsConf conf.TLSConf,
) CoreGrpcConnProvider {
	return &CoreGrpcConnProviderImpl{
		systemAccessTokenClient: systemAccessTokenClient,
		tlsConf:                 tlsConf,
	}
}
