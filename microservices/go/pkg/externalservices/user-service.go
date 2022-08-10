package externalservices

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb/userpbmapper"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"google.golang.org/grpc"
)

type ExtUserService interface {
	GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserDto]
	GetByAuthIdAsync(ctx context.Context, authId string) single.Single[userdtos.UserDto]
}

type ExtUserServiceImpl struct {
	grpcClientConf       conf.GrpcClientConf
	coreGrpcConnProvider CoreGrpcConnProvider
}

func (u ExtUserServiceImpl) GetByAuthIdAsync(ctx context.Context, authId string) single.Single[userdtos.UserDto] {
	return single.MapToAsync[userdtos.UserDto](ctx, u.GetByAuthId(ctx, authId))
}
func (u ExtUserServiceImpl) GetByAuthId(ctx context.Context, authId string) single.Single[userdtos.UserDto] {
	connectionSrc := u.coreGrpcConnProvider.CreateConnectionSingle(ctx, u.grpcClientConf.UserServiceAddress())
	userReplySrc := single.MapWithError(connectionSrc, func(conn *grpc.ClientConn) (*userpb.UserReply, error) {
		defer conn.Close()
		userService := userpb.NewUserServiceClient(conn)
		reply, err := userService.GetUserByAuthId(ctx, &userpb.AuthIdRequest{AuthId: authId})
		return reply, gtools.NewErrorResponseHandler(err).GetProcessedError()
	})
	return single.Map(userReplySrc, func(reply *userpb.UserReply) userdtos.UserDto {
		dto := userdtos.UserDto{}
		userpbmapper.MapUserReplyToUserDto(reply, &dto)
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
