package externalservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userkeypb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers/grpcmappers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/keydtos"
	"google.golang.org/grpc"
)

type ExtUserKeyService interface {
	GetKeyFromSession(
		ctx context.Context,
		userKeySessionDto keydtos.UserKeySessionDto,
	) single.Single[keydtos.UserKeyDto]
}

type ExtUserKeyServiceImpl struct {
	grpcClientConf       conf.GrpcClientConf
	coreGrpcConnProvider CoreGrpcConnProvider
}

func (e ExtUserKeyServiceImpl) GetKeyFromSession(
	ctx context.Context,
	userKeySessionDto keydtos.UserKeySessionDto,
) single.Single[keydtos.UserKeyDto] {
	connectionSrc := e.coreGrpcConnProvider.CreateConnectionSingle(ctx, e.grpcClientConf.KeyServiceAddress())
	replySrc := single.MapWithError(connectionSrc, func(conn *grpc.ClientConn) (*userkeypb.UserKey, error) {
		defer conn.Close()
		client := userkeypb.NewUserKeyServiceClient(conn)
		userKeySession := &userkeypb.UserKeySession{}
		grpcmappers.UserKeySessionDtoToUserKeySession(&userKeySessionDto, userKeySession)
		reply, err := client.GetKeyFromSession(ctx, userKeySession)
		return reply, gtools.NewErrorResponseHandler(err).GetProcessedError()
	})
	return single.Map(replySrc, func(res *userkeypb.UserKey) keydtos.UserKeyDto {
		dto := keydtos.UserKeyDto{}
		grpcmappers.UserKeyToUserKeyDto(res, &dto)
		return dto
	})
}

func NewExtUserKeyServiceImpl(
	grpcClientConf conf.GrpcClientConf,
	coreGrpcConnProvider CoreGrpcConnProvider,
) *ExtUserKeyServiceImpl {
	return &ExtUserKeyServiceImpl{grpcClientConf: grpcClientConf, coreGrpcConnProvider: coreGrpcConnProvider}
}
