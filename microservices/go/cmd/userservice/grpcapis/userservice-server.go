package grpcapis

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers/grpcmappers"
)

type UserServiceServerImpl struct {
	userpb.UnimplementedUserServiceServer
	userService services.UserService
}

func (u UserServiceServerImpl) GetUserById(ctx context.Context, request *userpb.IdRequest) (*userpb.UserReply, error) {
	userDto, err := u.userService.GetById(ctx, request.GetId())
	if err != nil {
		return nil, gtools.ProcessErrorToGrpcStatusError(gtools.ReadAction, err)
	}
	userReply := &userpb.UserReply{}
	grpcmappers.UserReadDtoToUserReply(&userDto, userReply)
	return userReply, err
}

func (u UserServiceServerImpl) GetUserByAuthId(
	ctx context.Context,
	request *userpb.AuthIdRequest,
) (*userpb.UserReply, error) {
	userDto, err := u.userService.GetByAuthId(ctx, request.GetAuthId())
	if err != nil {
		return nil, gtools.ProcessErrorToGrpcStatusError(gtools.ReadAction, err)
	}
	userReply := &userpb.UserReply{}
	grpcmappers.UserReadDtoToUserReply(&userDto, userReply)
	return userReply, err
}

func NewUserServiceServerImpl(userService services.UserService) *UserServiceServerImpl {
	return &UserServiceServerImpl{userService: userService}
}
