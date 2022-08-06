package grpcservers

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb/userpbmapper"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
)

type UserServiceServerImpl struct {
	userpb.UnimplementedUserServiceServer
	userService services.UserService
}

func (u UserServiceServerImpl) GetUserByAuthId(
	ctx context.Context,
	request *userpb.AuthIdRequest,
) (*userpb.UserReply, error) {
	userFindSrc := u.userService.GetByAuthId(ctx, request.GetAuthId())
	userReplySrc := single.Map(userFindSrc, func(userDto userdtos.UserDto) *userpb.UserReply {
		userReply := &userpb.UserReply{}
		userpbmapper.MapUserDtoToUserReply(&userDto, userReply)
		return userReply
	})
	res, err := single.RetrieveValue(ctx, userReplySrc)
	return res, gtools.ProcessErrorToGrpcStatusError(gtools.ReadAction, err)
}

func NewUserServiceServer(userService services.UserService) userpb.UserServiceServer {
	return &UserServiceServerImpl{userService: userService}
}
