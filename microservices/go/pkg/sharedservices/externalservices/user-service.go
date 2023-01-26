package externalservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers/grpcmappers"
	"google.golang.org/grpc"
)

type ExtUserService interface {
	GetByAuthId(ctx context.Context, authId string) (userdtos.UserReadDto, error)
	GetById(ctx context.Context, id string) (userdtos.UserReadDto, error)
}

type ExtUserServiceImpl struct {
	grpcClientConf       conf.GrpcClientConf
	coreGrpcConnProvider CoreGrpcConnProvider
}

func (u ExtUserServiceImpl) GetById(ctx context.Context, id string) (userDto userdtos.UserReadDto, err error) {
	conn, err := u.coreGrpcConnProvider.CreateConnectionSingle(ctx, u.grpcClientConf.UserServiceAddress())
	if err != nil {
		return userDto, err
	}
	defer func(conn *grpc.ClientConn) {
		if conErr := conn.Close(); conErr != nil {
			logger.Log.WithContext(ctx).WithError(conErr).Error()
		}
	}(conn)

	userService := userpb.NewUserServiceClient(conn)
	reply, err := userService.GetUserById(ctx, &userpb.IdRequest{Id: id})
	if err != nil {
		err = gtools.NewErrorResponseHandler(err).GetProcessedError()
		return userDto, err
	}

	dto := userdtos.UserReadDto{}
	grpcmappers.UserReplyToUserReadDto(reply, &dto)
	return dto, nil
}

func (u ExtUserServiceImpl) GetByAuthId(ctx context.Context, authId string) (userDto userdtos.UserReadDto, err error) {
	conn, err := u.coreGrpcConnProvider.CreateConnectionSingle(ctx, u.grpcClientConf.UserServiceAddress())
	if err != nil {
		return userDto, err
	}
	defer func(conn *grpc.ClientConn) {
		if conErr := conn.Close(); conErr != nil {
			logger.Log.WithContext(ctx).WithError(conErr).Error()
		}
	}(conn)

	userService := userpb.NewUserServiceClient(conn)
	reply, err := userService.GetUserByAuthId(ctx, &userpb.AuthIdRequest{AuthId: authId})
	if err != nil {
		err = gtools.NewErrorResponseHandler(err).GetProcessedError()
		return userDto, err
	}

	dto := userdtos.UserReadDto{}
	grpcmappers.UserReplyToUserReadDto(reply, &dto)
	return dto, nil
}

func NewExtUserServiceImpl(
	coreGrpcConnProvider CoreGrpcConnProvider,
	grpcClientConf conf.GrpcClientConf,
) *ExtUserServiceImpl {
	return &ExtUserServiceImpl{
		coreGrpcConnProvider: coreGrpcConnProvider,
		grpcClientConf:       grpcClientConf,
	}
}
