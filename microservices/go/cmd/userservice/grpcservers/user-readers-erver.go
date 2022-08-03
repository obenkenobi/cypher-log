package grpcservers

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/services"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/dtos/userdtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpb"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/grpc/userpbmapper"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
)

type UserReaderServerImpl struct {
	userpb.UnimplementedUserReaderServer
	userService services.UserService
}

func (u UserReaderServerImpl) GetUserByAuthId(ctx context.Context, request *userpb.AuthIdRequest) (*userpb.UserReply, error) {
	userFindSrc := u.userService.GetByAuthId(ctx, request.GetAuthId())
	userReplySrc := single.Map(userFindSrc, func(userDto userdtos.UserDto) *userpb.UserReply {
		userReply := &userpb.UserReply{}
		userpbmapper.MapUserDtoToUserReply(&userDto, userReply)
		return userReply
	})
	return single.RetrieveValue(ctx, userReplySrc)
}

func NewUserReaderServer(userService services.UserService) userpb.UserReaderServer {
	return &UserReaderServerImpl{userService: userService}
}
