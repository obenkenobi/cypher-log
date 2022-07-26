package grpcapis

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/gtools"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers/grpcmappers"
)

type UserServiceServerImpl struct {
	userpb.UnimplementedUserServiceServer
	userService services.UserService
}

func (u UserServiceServerImpl) GetUserById(ctx context.Context, request *userpb.IdRequest) (*userpb.UserReply, error) {
	userFindSrc := u.userService.GetById(ctx, request.GetId())
	userReplySrc := single.Map(userFindSrc, func(userDto userdtos.UserReadDto) *userpb.UserReply {
		userReply := &userpb.UserReply{}
		grpcmappers.UserReadDtoToUserReply(&userDto, userReply)
		return userReply
	})
	res, err := single.RetrieveValue(ctx, userReplySrc)
	return res, gtools.ProcessErrorToGrpcStatusError(gtools.ReadAction, err)
}

func (u UserServiceServerImpl) GetUserByAuthId(
	ctx context.Context,
	request *userpb.AuthIdRequest,
) (*userpb.UserReply, error) {
	userFindSrc := u.userService.GetByAuthId(ctx, request.GetAuthId())
	userReplySrc := single.Map(userFindSrc, func(userDto userdtos.UserReadDto) *userpb.UserReply {
		userReply := &userpb.UserReply{}
		grpcmappers.UserReadDtoToUserReply(&userDto, userReply)
		return userReply
	})
	res, err := single.RetrieveValue(ctx, userReplySrc)
	return res, gtools.ProcessErrorToGrpcStatusError(gtools.ReadAction, err)
}

func NewUserServiceServerImpl(userService services.UserService) *UserServiceServerImpl {
	return &UserServiceServerImpl{userService: userService}
}
