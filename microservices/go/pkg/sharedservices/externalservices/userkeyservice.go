package externalservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userkeypb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers/grpcmappers"
	"google.golang.org/grpc"
)

type ExtUserKeyService interface {
	GetKeyFromSession(
		ctx context.Context,
		userKeySessionDto commondtos.UKeySessionDto,
	) (keydtos.UserKeyDto, error)
}

type ExtUserKeyServiceImpl struct {
	grpcClientConf       conf.GrpcClientConf
	coreGrpcConnProvider CoreGrpcConnProvider
}

func (e ExtUserKeyServiceImpl) GetKeyFromSession(
	ctx context.Context,
	userKeySessionDto commondtos.UKeySessionDto,
) (keyDto keydtos.UserKeyDto, err error) {
	conn, err := e.coreGrpcConnProvider.CreateConnectionSingle(ctx, e.grpcClientConf.KeyServiceAddress())
	if err != nil {
		return keyDto, err
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
	}(conn)

	client := userkeypb.NewUserKeyServiceClient(conn)

	userKeySession := &userkeypb.UserKeySession{}
	grpcmappers.UserKeySessionDtoToUserKeySession(&userKeySessionDto, userKeySession)

	reply, err := client.GetKeyFromSession(ctx, userKeySession)
	if err != nil {
		err = gtools.NewErrorResponseHandler(err).GetProcessedError()
		return keyDto, err
	}

	dto := keydtos.UserKeyDto{}
	grpcmappers.UserKeyToUserKeyDto(reply, &dto)
	return dto, nil
}

func NewExtUserKeyServiceImpl(
	grpcClientConf conf.GrpcClientConf,
	coreGrpcConnProvider CoreGrpcConnProvider,
) *ExtUserKeyServiceImpl {
	return &ExtUserKeyServiceImpl{grpcClientConf: grpcClientConf, coreGrpcConnProvider: coreGrpcConnProvider}
}
