package externalservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers/grpcmappers"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/userdtos"
	"google.golang.org/grpc"
)

type ExtUserService interface {
	GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserReadDto]
	GetById(ctx context.Context, id string) single.Single[userdtos.UserReadDto]
}

type ExtUserServiceImpl struct {
	grpcClientConf       conf.GrpcClientConf
	coreGrpcConnProvider CoreGrpcConnProvider
}

func (u ExtUserServiceImpl) GetById(ctx context.Context, id string) single.Single[userdtos.UserReadDto] {
	connectionSrc := u.coreGrpcConnProvider.CreateConnectionSingle(ctx, u.grpcClientConf.UserServiceAddress())
	userReplySrc := single.MapWithError(
		connectionSrc,
		func(conn *grpc.ClientConn) (reply *userpb.UserReply, err error) {
			defer func(conn *grpc.ClientConn) { err = conn.Close() }(conn)
			userService := userpb.NewUserServiceClient(conn)
			reply, err = userService.GetUserById(ctx, &userpb.IdRequest{Id: id})
			return reply, gtools.NewErrorResponseHandler(err).GetProcessedError()
		},
	)
	return single.Map(userReplySrc, func(reply *userpb.UserReply) userdtos.UserReadDto {
		dto := userdtos.UserReadDto{}
		grpcmappers.MapUserReplyToUserReadDto(reply, &dto)
		return dto
	})
}

func (u ExtUserServiceImpl) GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserReadDto] {
	connectionSrc := u.coreGrpcConnProvider.CreateConnectionSingle(ctx, u.grpcClientConf.UserServiceAddress())
	userReplySrc := single.MapWithError(connectionSrc, func(conn *grpc.ClientConn) (*userpb.UserReply, error) {
		defer conn.Close()
		userService := userpb.NewUserServiceClient(conn)
		reply, err := userService.GetUserByAuthId(ctx, &userpb.AuthIdRequest{AuthId: authId})
		return reply, gtools.NewErrorResponseHandler(err).GetProcessedError()
	})
	return single.Map(userReplySrc, func(reply *userpb.UserReply) userdtos.UserReadDto {
		dto := userdtos.UserReadDto{}
		grpcmappers.MapUserReplyToUserReadDto(reply, &dto)
		return dto
	})
}

func NewExtUserService(
	coreGrpcConnProvider CoreGrpcConnProvider,
	grpcClientConf conf.GrpcClientConf,
) ExtUserService {
	return &ExtUserServiceImpl{
		coreGrpcConnProvider: coreGrpcConnProvider,
		grpcClientConf:       grpcClientConf,
	}
}
